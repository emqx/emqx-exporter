#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from collections import OrderedDict
from metrics import subcharts
from submetrics import subchart_metrics
from common import (
    create_nodes_template_list,
    create_dashboard,
    create_panel,
    create_row,
    create_timeseries,
    generate_timeseries_overrides,
)

key = subcharts["messages"]
title = subchart_metrics[key]['title']
uid = subchart_metrics[key]['uid']
panels = OrderedDict(subchart_metrics[key]['panels'])


def emqx_messages_dashboard(dashboard):
    dashboard.add_row(create_row(
        title=title,
        collapsed=False))

    for _panel_key, panel_value in panels.items():
        create_panel(
            dashboard,
            create_timeseries(
                title=panel_value['title'],
                span=4,
                overrides=generate_timeseries_overrides(
                    panel_value['targets']),
                legendDisplayMode="table",
                legendCalcs=[
                    "lastNotNull",
                    "min",
                    "max",
                    "mean",
                    "sum",
                ]),
            panel_value['targets'])


if __name__ == '__main__':
    # create a dashboard
    dashboard = create_dashboard(
        title=title,
        uid=uid,
        templating=create_nodes_template_list(),
    )

    emqx_messages_dashboard(dashboard)

    dashboard.auto_panel_ids()
    print(dashboard.generate_dashboard())
