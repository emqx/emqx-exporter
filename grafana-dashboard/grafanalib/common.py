import io
import re
import grafanalib.core as G
from grafanalib._gen import write_dashboard
from submetrics import subchart_metrics


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


class EmqxDashboard(G.Dashboard):
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


def find_panel_info(subchart, key):
    # Check if the specified subchart exists and is a dictionary
    if subchart in subchart_metrics and isinstance(subchart_metrics[subchart], dict):
        subchart_data = subchart_metrics[subchart]
        uid = subchart_data.get('uid')

        # Check if 'panels' key exists and is a dictionary in the subchart
        if 'panels' in subchart_data and isinstance(subchart_data['panels'], dict):
            panels = subchart_data['panels']

            # Iterate over panels
            for index, (panel_key, panel_value) in enumerate(panels.items()):
                if panel_key == key:
                    title = panel_value.get('title')
                    return uid, title, index

    return None, None, -1


def generate_a_url(uid, title, index):
    return f"/d/{uid}/{title}?orgId=1&refresh=10s&var-node=All&viewPanel={index}"


def format_string(s):
    # Remove 'subchart_' prefix and replace underscores with spaces
    return ' '.join(word.capitalize() for word in s.replace('subchart_', '').split('_'))


def create_data_links(subchart_links):
    data_links = []

    subchart, links = subchart_links[0], subchart_links[1]
    for link in links:
        uid, title, index = find_panel_info(subchart, link)
        if title is not None:
            data_links.append(G.DataLink(
                title=f"Show {format_string(link)} Node Detail",
                linkUrl=generate_a_url(uid, title, index+1),
                isNewTab=True))
        else:
            raise ValueError(f"can not find subchart link '{link}'")

    return data_links


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
            dataSource="${datasource}",
            query="label_values(up, cluster)",
            includeAll=False,
            multi=False,
            options=[],
            refresh=1,
            regex="",
            hide=0,
        ),
    ])


def create_nodes_template_list():
    return G.Templating([
        G.Template(
            hide=0,
            includeAll=False,
            label="datasource",
            multi=False,
            name="datasource",
            options=[],
            query="prometheus",
            refresh=1,
            regex="",
            type="datasource"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=False,
            multi=False,
            name="cluster",
            options=[],
            query={
                "query": "label_values(up, cluster)",
                "refId": "StandardVariableQuery"
            },
            refresh=1,
            regex="",
            sort=2,
            type="query"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=True,
            multi=True,
            name="node",
            options=[],
            query={
                "query": "label_values(up{from=\"emqx\",cluster=\"$cluster\"}, instance)",
                "refId": "StandardVariableQuery"
            },
            refresh=1,
            regex="",
            sort=1,
            type="query"
        )
    ])


def create_auth_template_list():
    return G.Templating([
        G.Template(
            hide=0,
            includeAll=False,
            label="datasource",
            multi=False,
            name="datasource",
            options=[],
            query="prometheus",
            refresh=1,
            regex="",
            type="datasource"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=False,
            multi=False,
            name="cluster",
            options=[],
            query={
                "query": "label_values(up, cluster)",
                "refId": "StandardVariableQuery"
            },
            refresh=1,
            regex="",
            sort=2,
            type="query"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=True,
            multi=True,
            name="node",
            options=[],
            query={
                "query": "label_values({cluster=\"$cluster\", from=\"exporter\"}, node)",
                "refId": "StandardVariableQuery"
            },
            refresh=2,
            regex="",
            sort=1,
            type="query"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=True,
            multi=True,
            name="resource",
            options=[],
            query={
                "query": "label_values(emqx_authentication_resource_status{cluster=\"$cluster\"}, resource)",
                "refId": "StandardVariableQuery"
            },
            refresh=1,
            regex="",
            sort=1,
            type="query"
        )
    ])


def create_rule_engine_template_list():
    return G.Templating([
        G.Template(
            hide=0,
            includeAll=False,
            label="datasource",
            multi=False,
            name="datasource",
            options=[],
            query="prometheus",
            refresh=1,
            regex="",
            type="datasource"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=False,
            multi=False,
            name="cluster",
            options=[],
            query={
                "query": "label_values(up, cluster)",
                "refId": "StandardVariableQuery"
            },
            refresh=1,
            regex="",
            sort=2,
            type="query"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=True,
            multi=True,
            name="node",
            options=[],
            query={
                "query": "label_values({cluster=\"$cluster\", from=\"exporter\"}, node)",
                "refId": "StandardVariableQuery"
            },
            refresh=2,
            regex="",
            sort=1,
            type="query"
        ),
        G.Template(
            dataSource={
                "type": "prometheus",
                "uid": "${datasource}"
            },
            hide=0,
            includeAll=True,
            multi=True,
            name="rule",
            options=[],
            query={
                "query": "label_values({cluster=\"$cluster\", node=~\"$node\", from=\"exporter\"}, rule)",
                "refId": "StandardVariableQuery"
            },
            refresh=1,
            regex="",
            sort=1,
            type="query"
        )
    ])


def create_dashboard(**kwargs):
    default_kwargs = {
        "title": "emqx-overview",
        "time": create_time("now-5m", "now"),
        "refresh": "5s",
        "tags": ["EMQX", "MQTT"],
        "templating": create_template_list(),
    }

    merged_kwargs = {**default_kwargs, **kwargs}

    return EmqxDashboard(**merged_kwargs)


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
        "dataSource": "${datasource}",
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
        "dataSource": "${datasource}",
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
        "dataSource": "${datasource}",
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
        "dataSource": "${datasource}",
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
