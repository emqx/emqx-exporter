#!/usr/bin/env python3
# -*- coding: utf-8 -*-


import argparse

from common import (
    create_timeseries,
    create_row,
    create_panel,
    create_dashboard,
    create_stat,
    create_gauge,
    create_table,
    create_gridpos,
    create_data_links,

    generate_transformations,
    generate_table_overrides,
    generate_timeseries_overrides,
    thresholds_2_steps,
    thresholds_3_steps,
)
from metrics import metrics


def emqx_dashboard(dashboard, is_ee, version):

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

    if is_ee:
        create_panel(
            dashboard,
            create_table(
                title=metrics['license']['title'],
                gridPos=create_gridpos(metrics['license']['gridPos']),
                transformations=generate_transformations(
                    metrics['license']['targets']),
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
            overrides=generate_timeseries_overrides(
                metrics['nodes_running']['targets']),
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
            links=create_data_links(
                metrics['rule_engine_last_5m_exec_rate']['subchart_links']),
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

    if version == 5:
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
            links=create_data_links(
                metrics['rule_engine_current_exec_rate']['subchart_links']),
            span=3),
        metrics['rule_engine_current_exec_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_exec_success']['title'],
            links=create_data_links(
                metrics['rule_engine_exec_success']['subchart_links']),
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_exec_success']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_exec_failure']['title'],
            links=create_data_links(
                metrics['rule_engine_exec_failure']['subchart_links']),
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_exec_failure']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_action_success']['title'],
            links=create_data_links(
                metrics['rule_engine_action_success']['subchart_links']),
            span=3,
            thresholds=thresholds_2_steps),
        metrics['rule_engine_action_success']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['rule_engine_action_failure']['title'],
            links=create_data_links(
                metrics['rule_engine_action_failure']['subchart_links']),
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
            links=create_data_links(
                metrics['client_connection_events']['subchart_links']),
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
            links=create_data_links(
                metrics['client_sub_events']['subchart_links']),
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
            links=create_data_links(
                metrics['client_connect_auth_events']['subchart_links']),
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

    if version == 5:
        create_panel(
            dashboard,
            create_timeseries(
                title=metrics['client_acl_auth_events_v5']['title'],
                links=create_data_links(
                    metrics['client_acl_auth_events_v5']['subchart_links']),
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
            metrics['client_acl_auth_events_v5']['targets'])
    else:
        create_panel(
            dashboard,
            create_timeseries(
                title=metrics['client_acl_auth_events_v4']['title'],
                links=create_data_links(
                    metrics['client_acl_auth_events_v4']['subchart_links']),
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
            metrics['client_acl_auth_events_v4']['targets'])

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
            overrides=generate_timeseries_overrides(
                metrics['packets_connections']['targets']),
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
            overrides=generate_timeseries_overrides(
                metrics['packets_disconnections']['targets']),
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
            overrides=generate_timeseries_overrides(
                metrics['packets_publish']['targets']),
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
            overrides=generate_timeseries_overrides(
                metrics['packets_subscribe_and_unsubscribe']['targets']),
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
            links=create_data_links(
                metrics['messages_count']['subchart_links']),
            span=4,
            overrides=generate_timeseries_overrides(
                metrics['messages_count']['targets']),
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
            overrides=generate_timeseries_overrides(
                metrics['messages_qos_received']['targets']),
            thresholds=thresholds_2_steps),
        metrics['messages_qos_received']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=metrics['cluster_traffic_statistics']['title'],
            span=4,
            unit="decbytes",
            overrides=generate_timeseries_overrides(
                metrics['cluster_traffic_statistics']['targets']),
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
            transformations=generate_transformations(
                metrics['data_bridge_status']['targets']),
            overrides=generate_table_overrides(metrics['data_bridge_status']['targets'])),
        metrics['data_bridge_status']['targets'],
        instant=True,
        format=metrics['data_bridge_status']['format'])

    create_panel(
        dashboard,
        create_table(
            title=metrics['rule_engine_execute_count']['title'],
            span=8,
            transformations=generate_transformations(
                metrics['rule_engine_execute_count']['targets']),
            overrides=generate_table_overrides(metrics['rule_engine_execute_count']['targets'])),
        metrics['rule_engine_execute_count']['targets'],
        instant=True,
        format=metrics['rule_engine_execute_count']['format'])

    if version == 5:
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
                links=create_data_links(
                    metrics['authenticate_count']['subchart_links']),
                span=4,
                transformations=generate_transformations(
                    metrics['authenticate_count']['targets']),
                overrides=generate_table_overrides(metrics['authenticate_count']['targets'])),
            metrics['authenticate_count']['targets'],
            instant=True,
            format=metrics['authenticate_count']['format'])

        create_panel(
            dashboard,
            create_timeseries(
                title=metrics['authenticate_current_exec_rate']['title'],
                links=create_data_links(
                    metrics['authenticate_current_exec_rate']['subchart_links']),
                span=4,
            ),
            metrics['authenticate_current_exec_rate']['targets'])

        create_panel(
            dashboard,
            create_timeseries(
                title=metrics['authenticate_last_5m_exec_rate']['title'],
                links=create_data_links(
                    metrics['authenticate_last_5m_exec_rate']['subchart_links']),
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
                links=create_data_links(
                    metrics['authorize_count']['subchart_links']),
                span=4,
                transformations=generate_transformations(
                    metrics['authorize_count']['targets']),
                overrides=generate_table_overrides(metrics['authorize_count']['targets'])),
            metrics['authorize_count']['targets'],
            instant=True,
            format=metrics['authorize_count']['format'])

        create_panel(
            dashboard,
            create_timeseries(
                title=metrics['authorize_current_exec_rate']['title'],
                links=create_data_links(
                    metrics['authorize_current_exec_rate']['subchart_links']),
                span=4),
            metrics['authorize_current_exec_rate']['targets'])

        create_panel(
            dashboard,
            create_timeseries(
                title=metrics['authorize_last_5m_exec_rate']['title'],
                links=create_data_links(
                    metrics['authorize_last_5m_exec_rate']['subchart_links']),
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


if __name__ == '__main__':
    is_ee = True
    ver = 5

    parser = argparse.ArgumentParser(description="Argument Parser")

    # Define a single optional argument for category (-c or --category)
    parser.add_argument(
        "-e", "--edition", choices=["ee", "ce"], help="Set EMQX Enterprise(ee) or EMQX Community(ce)")
    parser.add_argument("-v", "--version", type=int,
                        choices=[4, 5], help="Set EMQX Version")

    args = parser.parse_args()

    if args.edition == "ce":
        is_ee = False

    if args.version == 4:
        ver = 4

    # create a dashboarnkd
    dashboard = create_dashboard(title="Overview", uid="overview")

    emqx_dashboard(dashboard, is_ee, ver)

    dashboard.auto_panel_ids()
    print(dashboard.generate_dashboard())
