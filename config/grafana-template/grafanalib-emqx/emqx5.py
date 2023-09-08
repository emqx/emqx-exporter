#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
NAME:
    emqx5.py

DESCRIPTION:
    This script creates Grafana dashboards using Grafanalib, and a static table
    which defines metrics/dashboards.

    The resulting dashboard can be easily uploaded to Grafana with associated script:

        upload_grafana_dashboard.sh

USAGE:
    Create and upload the dashboard:

    ./emqx5.py > dash.json
    ./upload_grafana_dashboard.sh dash.json

"""

import io
import re
import grafanalib.core as G
from grafanalib._gen import write_dashboard

from metrics import metrics

thresholds_2_steps = [
    {
        "color": "green",
        "value": None
    },
    {
        "color": "red",
        "value": 80
    }
]

thresholds_3_steps = [
    {
        "color": "green",
        "value": None
    },
    {
        "color": "orange",
        "value": 80
    },
    {
        "color": "red",
        "value": 90
    }
]


class Emqx5Dashboard(G.Dashboard):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.targets = []

    def add_row(self, row):
        self.rows.append(row)

    def add_panel(self, panel):
        self.rows[-1].panels.append(panel)

    def auto_panel_ids(self):
        """Auto-number panels in dashboard"""
        ref_id = 1
        for row in self.rows:
            for panel in row.panels:
                panel.id = ref_id
                ref_id = ref_id + 1
        return self

    def generate_dashboard(self):
        s = io.StringIO()
        write_dashboard(self, s)
        return s.getvalue()


class TimeSeries(G.TimeSeries):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.targets = []

    def add_target(self, target):
        self.targets.append(target)


class Stat(G.Stat):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.targets = []

    def add_target(self, target):
        self.targets.append(target)


class Gauge(G.GaugePanel):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.targets = []

    def add_target(self, target):
        self.targets.append(target)


class Table(G.Table):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.targets = []

    def add_target(self, target):
        self.targets.append(target)


# Helper functions to create dashboard components

def create_time(start, end):
    return G.Time(start=start, end=end)


def create_template_list():
    return G.Templating([
        G.Template(
            name="datasource",
            label="datasource",
            query="prometheus",
            type="datasource",
            includeAll=False,
            multi=False,
            options=[],
            refresh=1,
            regex="",
            hide=0,
        ),
        G.Template(
            name="cluster",
            type="query",
            dataSource="$datasource",
            query="label_values(up, cluster)",
            includeAll=False,
            multi=False,
            options=[],
            refresh=1,
            regex="",
            hide=0,
        ),
    ])


def create_dashboard(**kwargs):
    default_kwargs = {
        "title": "emqx-test",
        "time": create_time("now-5m", "now"),
        "refresh": "5s",
        "tags": ["test"],
        "templating": create_template_list(),
    }

    merged_kwargs = {**default_kwargs, **kwargs}

    return Emqx5Dashboard(**merged_kwargs)


def create_row(title, collapsed=True):
    return G.Row(title=title, collapse=collapsed)


def create_target(target, instant, format):
    return G.Target(
        expr=target['expr'],
        legendFormat=target['legendFormat'],
        intervalFactor=1,
        instant=instant,
        format=format,
        refId='_'.join(target['legendFormat'].lower().split()))


def create_gridpos(pos):
    return G.GridPos(
        h=pos['h'],
        w=pos['w'],
        x=pos['x'],
        y=pos['y']
    )


def create_timeseries(**kwargs):
    default_kwargs = {
        "dataSource": "prometheus",
        "gridPos": create_gridpos({"h": 8, "w": 6, "x": 0, "y": 14}),
        "span": 2,
        "scaleDistributionType": "linear",
        "tooltipMode": "multi",
        "lineInterpolation": "linear",
        "showPoints": "never",
        "gradientMode": "opacity",
        "legendDisplayMode": "list",
        "legendPlacement": "bottom",
        "unit": "short",
    }

    merged_kwargs = {**default_kwargs, **kwargs}
    return TimeSeries(**merged_kwargs)


def create_stat(**kwargs):
    default_kwargs = {
        "dataSource": "prometheus",
        "gridPos": create_gridpos({"h": 3, "w": 6, "x": 0, "y": 1}),
        "span": "2",
        "mappings": [
            {
                "options": {
                    "0": {
                        "color": "red",
                        "index": 2,
                        "text": "Unknown"
                    },
                    "1": {
                        "color": "red",
                        "index": 1,
                        "text": "Unhealthy"
                    },
                    "2": {
                        "color": "green",
                        "index": 0,
                        "text": "Healthy"
                    }
                },
                "type": "value"
            },
        ],

        "reduceCalc": "lastNotNull",
        "fields": "/^Status$/",
        "orientation": "auto",
        "textMode": "auto",
        "colorMode": "background",
        "graphMode": "none",
        "alignment": "auto",
    }

    merged_kwargs = {**default_kwargs, **kwargs}
    return Stat(**merged_kwargs)


def create_gauge(**kwargs):
    default_kwargs = {
        "dataSource": "prometheus",
        "gridPos": create_gridpos({"h": 6, "w": 4, "x": 6, "y": 1}),
        "span": 2,
        "max": None,
        "thresholdType": "percentage",
        "thresholds": thresholds_3_steps,
    }

    merged_kwargs = {**default_kwargs, **kwargs}
    return Gauge(**merged_kwargs)


def create_table(**kwargs):
    default_kwargs = {
        "dataSource": "prometheus",
        "gridPos": create_gridpos({"h": 6, "w": 10, "x": 0, "y": 42}),
        "span": 2,
        "thresholds": thresholds_2_steps,
        # "overrides": [
        # ],
        "showHeader": True,
        "filterable": True,
    }

    merged_kwargs = {**default_kwargs, **kwargs}
    return Table(**merged_kwargs)


def create_panel(dashboard, panel, targets, instant=False, format="timeseries"):
    for target in targets:
        panel.add_target(create_target(target, instant, format))

    dashboard.add_panel(panel)


def extract_fields_from_expr(expr):
    pattern = r'\((.*?)\)'

    matches = re.search(pattern, expr)
    if matches:
        content = matches.group(1)
        return [x.strip() for x in content.split(',')]
    else:
        raise Exception("No match found.")


# Generate transformations with LegendFormat of targets
def generate_transformations(targets):
    transformations = [
        {
            "id": "merge",
            "options": {}
        },
        {
            "id": "filterFieldsByName",
            "options": {
                "include": {
                    "names": [
                    ]
                }
            }
        }
    ]
    fields = extract_fields_from_expr(targets[0]['expr'])

    for field in fields:
        transformations[1]["options"]["include"]["names"].append(field)

    for target in targets:
        transformations[1]["options"]["include"]["names"].append(
            "Value #" + "_".join(target['legendFormat'].lower().split())
        )
    return transformations


def generate_timeseries_overrides(targets):
    overrides = []
    for target in targets:
        item = {}
        if 'color' in target:
            item = {
                "matcher": {
                    "id": "byName",
                    "options": target['legendFormat']
                },
                "properties": [
                    {
                        "id": "color",
                        "value": {
                            "fixedColor": target['color'],
                            "mode": "fixed"
                        }
                    },
                ]
            }
            overrides.append(item)
    return overrides


# Generate overrides with LegendFormat of targets
def generate_table_overrides(targets):
    overrides = []
    for target in targets:
        item = {}
        if 'thresholds' in target:
            item = {
                "matcher": {
                    "id": "byName",
                    "options": "Value #" + "_".join(target['legendFormat'].lower().split())
                },
                "properties": [
                    {
                        "id": "displayName",
                        "value": target['legendFormat'],
                    },
                    {
                        "id": "thresholds",
                        "value": target['thresholds'],
                    },
                    {
                        "id": "custom.cellOptions",
                        "value": {
                            "mode": "gradient",
                            "type": "color-background"
                        }
                    },
                ]
            }
        elif 'mappings' in target:
            item = {
                "matcher": {
                    "id": "byName",
                    "options": "Value #" + "_".join(target['legendFormat'].lower().split())
                },
                "properties": [
                    {
                        "id": "displayName",
                        "value": target['legendFormat'],
                    },
                    {
                        "id": "mappings",
                        "value": target['mappings'],
                    },
                    {
                        "id": "custom.cellOptions",
                        "value": {
                            "mode": "gradient",
                            "type": "color-background"
                        }
                    },
                ]
            }
        elif 'datetime' in target:
            item = {
                "matcher": {
                    "id": "byName",
                    "options": "Value #" + "_".join(target['legendFormat'].lower().split())
                },
                "properties": [
                    {
                        "id": "unit",
                        "value": "dateTimeAsLocal"
                    },
                ]
            }
        else:
            item = {
                "matcher": {
                    "id": "byName",
                    "options": "Value #" + "_".join(target['legendFormat'].lower().split())
                },
                "properties": [
                    {
                        "id": "displayName",
                        "value": target['legendFormat']
                    }
                ]
            }

        overrides.append(item)
    return overrides


if __name__ == '__main__':

    # create a dashboard
    dashboard = create_dashboard()

    ##############################
    # General
    ##############################

    dashboard.add_row(create_row(
        title=metrics['general_row']['title'],
        collapsed=False))

    create_panel(
        dashboard,
        create_stat(title=metrics['cluster_status']['title']),
        metrics['cluster_status']['targets'],
        format=metrics['cluster_status']['format'])

    create_panel(
        dashboard,
        create_table(
            title=metrics['license']['title'],
            gridPos=create_gridpos(metrics['license']['gridPos']),
            transformations=generate_transformations(metrics['license']['targets']),
            overrides=generate_table_overrides(metrics['license']['targets'])),
        metrics['license']['targets'],
        instant=True,
        format=metrics['license']['format'])

    create_panel(
        dashboard,
        create_gauge(title=metrics['active_connections']['title']),
        metrics['active_connections']['targets'],
        format=metrics['active_connections']['format'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['cluster_message_rate']['title'],
            thresholds=thresholds_2_steps,
            gridPos=create_gridpos(metrics['cluster_message_rate']['gridPos'])),
        metrics['cluster_message_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['nodes_running']['title'],
            gradientMode="none",
            fillOpacity=15,
            thresholds=thresholds_2_steps,
            overrides=generate_timeseries_overrides(metrics['nodes_running']['targets']),
            gridPos=create_gridpos(metrics['nodes_running']['gridPos'])),
        metrics['nodes_running']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['exporter_latency']['title'], unit='s',
            gradientMode="none",
            fillOpacity=15,
            thresholds=thresholds_2_steps,
            gridPos=create_gridpos(metrics['exporter_latency']['gridPos'])),
        metrics['exporter_latency']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['sessions']['title'],
            gradientMode="none",
            fillOpacity=15,
            thresholdsStyleMode="area",
            span=3,
            thresholdType="percentage",
            thresholds=thresholds_3_steps),
        metrics['sessions']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['connections']['title'],
            span=3),
        metrics['connections']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['subscriptions']['title'],
            span=3),
        metrics['subscriptions']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_last_5m_exec_rate']['title'],
            span=3),
        metrics['rule_engine_last_5m_exec_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['connect_auth']['title'],
            span=3,
            thresholds=thresholds_2_steps),
        metrics['connect_auth']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['acl_auth']['title'],
            span=3,
            thresholds=thresholds_2_steps),
        metrics['acl_auth']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['data_bridge_queuing']['title'],
            span=3),
        metrics['data_bridge_queuing']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_current_exec_rate']['title'],
            span=3),
        metrics['rule_engine_current_exec_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_exec_success']['title'],
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_exec_success']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_exec_failure']['title'],
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_exec_failure']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_action_success']['title'],
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_action_success']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_action_failure']['title'],
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_action_failure']['targets'])

    ##############################
    # Event
    ##############################

    dashboard.add_row(create_row(
        title=metrics['event_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['client_connection_events']['title'],
            span=3,
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['client_connection_events']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['client_sub_events']['title'],
            span=3,
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['client_sub_events']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['client_connect_auth_events']['title'],
            span=3,
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['client_connect_auth_events']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['client_acl_auth_events']['title'],
            span=3,
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['client_acl_auth_events']['targets'])

    ##############################
    # Packets
    ##############################

    dashboard.add_row(create_row(
        title=metrics['packets_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['packets_connections']['title'],
            span=3,
            overrides=generate_timeseries_overrides(metrics['packets_connections']['targets']),
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['packets_connections']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['packets_disconnections']['title'],
            span=3,
            overrides=generate_timeseries_overrides(metrics['packets_disconnections']['targets']),
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['packets_disconnections']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['packets_publish']['title'],
            span=3,
            overrides=generate_timeseries_overrides(metrics['packets_publish']['targets']),
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['packets_publish']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['packets_subscribe_and_unsubscribe']['title'],
            span=3,
            overrides=generate_timeseries_overrides(metrics['packets_subscribe_and_unsubscribe']['targets']),
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ],
            thresholds=thresholds_2_steps),
        metrics['packets_subscribe_and_unsubscribe']['targets'])

    ##############################
    # Messages
    ##############################

    dashboard.add_row(create_row(
        title=metrics['messages_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['messages_count']['title'],
            span=4,
            overrides=generate_timeseries_overrides(metrics['messages_count']['targets']),
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ]),
        metrics['messages_count']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['messages_qos_received']['title'],
            span=4,
            overrides=generate_timeseries_overrides(metrics['messages_qos_received']['targets']),
            thresholds=thresholds_2_steps),
        metrics['messages_qos_received']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['cluster_traffic_statistics']['title'],
            span=4,
            unit="decbytes",
            overrides=generate_timeseries_overrides(metrics['cluster_traffic_statistics']['targets']),
            legendDisplayMode="table",
            legendCalcs=[
                "lastNotNull",
                "min",
                "max",
                "mean",
                "sum",
            ]),
        metrics['cluster_traffic_statistics']['targets'])

    ##############################
    # Rule Engine
    ##############################

    dashboard.add_row(create_row(
        title=metrics['rule_engine_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_table(
            title=metrics['data_bridge_status']['title'],
            span=4,
            transformations=generate_transformations(metrics['data_bridge_status']['targets']),
            overrides=generate_table_overrides(metrics['data_bridge_status']['targets'])),
        metrics['data_bridge_status']['targets'],
        instant=True,
        format=metrics['data_bridge_status']['format'])

    create_panel(
        dashboard,
        create_table(
            title=metrics['rule_engine_execute_count']['title'],
            span=8,
            transformations=generate_transformations(metrics['rule_engine_execute_count']['targets']),
            overrides=generate_table_overrides(metrics['rule_engine_execute_count']['targets'])),
        metrics['rule_engine_execute_count']['targets'],
        instant=True,
        format=metrics['rule_engine_execute_count']['format'])

    ##############################
    # Connect Auth
    ##############################

    dashboard.add_row(create_row(
        title=metrics['connect_auth_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_table(
            title=metrics['authenticate_count']['title'],
            span=4,
            transformations=generate_transformations(metrics['authenticate_count']['targets']),
            overrides=generate_table_overrides(metrics['authenticate_count']['targets'])),
        metrics['authenticate_count']['targets'],
        instant=True,
        format=metrics['authenticate_count']['format'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['authenticate_current_exec_rate']['title'],
            span=4,
        ),
        metrics['authenticate_current_exec_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['authenticate_last_5m_exec_rate']['title'],
            span=4
        ),
        metrics['authenticate_last_5m_exec_rate']['targets'])

    ##############################
    # ACL Auth
    ##############################

    dashboard.add_row(create_row(
        title=metrics['acl_auth_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_table(
            title=metrics['authorize_count']['title'],
            span=4,
            transformations=generate_transformations(metrics['authorize_count']['targets']),
            overrides=generate_table_overrides(metrics['authorize_count']['targets'])),
        metrics['authorize_count']['targets'],
        instant=True,
        format=metrics['authorize_count']['format'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['authorize_current_exec_rate']['title'],
            span=4),
        metrics['authorize_current_exec_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['authorize_last_5m_exec_rate']['title'],
            span=4),
        metrics['authorize_last_5m_exec_rate']['targets'])

    ##############################
    # System
    ##############################

    dashboard.add_row(create_row(
        title=metrics['system_row']['title'],
        collapsed=True))

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['last_1m_cpu_load']['title'],
            span=4,
            thresholds=thresholds_2_steps),
        metrics['last_1m_cpu_load']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['last_5m_cpu_load']['title'],
            span=4,
            thresholds=thresholds_2_steps),
        metrics['last_5m_cpu_load']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['last_15m_cpu_load']['title'],
            span=4,
            thresholds=thresholds_2_steps),
        metrics['last_15m_cpu_load']['targets'])

    dashboard.auto_panel_ids()
    print(dashboard.generate_dashboard())
