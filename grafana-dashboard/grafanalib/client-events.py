#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse

from collections import OrderedDict

from metrics import subcharts
from submetrics import subchart_metrics
from common import (
    create_nodes_template_list,
    create_dashboard,
    create_panel,
    create_row,
    create_timeseries,
    thresholds_2_steps,
)


key = subcharts["client-events"]
title = subchart_metrics[key]['title']
uid = subchart_metrics[key]['uid']
panels = OrderedDict(subchart_metrics[key]['panels'])


def emqx_client_events_dashboard(dashboard, version):
    dashboard.add_row(create_row(
        title=title,
        collapsed=False))

    for panel_key, panel_value in panels.items():
        if panel_key.startswith('subchart_client_acl_auth_events'):
            continue

        create_panel(
            dashboard,
            create_timeseries(
                title=panel_value['title'],
                span=6,
                legendDisplayMode="table",
                legendCalcs=[
                    "lastNotNull",
                    "min",
                    "max",
                    "mean",
                    "sum",
                ],
                thresholds=thresholds_2_steps),
            panel_value['targets'])

    if version == 5:
        create_panel(
            dashboard,
            create_timeseries(
                title=panels['subchart_client_acl_auth_events_v5']['title'],
                span=6,
                legendDisplayMode="table",
                legendCalcs=[
                    "lastNotNull",
                    "min",
                    "max",
                    "mean",
                    "sum",
                ],
                thresholds=thresholds_2_steps),
            panels['subchart_client_acl_auth_events_v5']['targets'])
    else:
        create_panel(
            dashboard,
            create_timeseries(
                title=panels['subchart_client_acl_auth_events_v4']['title'],
                span=6,
                legendDisplayMode="table",
                legendCalcs=[
                    "lastNotNull",
                    "min",
                    "max",
                    "mean",
                    "sum",
                ],
                thresholds=thresholds_2_steps),
            panels['subchart_client_acl_auth_events_v4']['targets'])


if __name__ == '__main__':
    ver = 5

    parser = argparse.ArgumentParser(description="Argument Parser")

    # Define a single optional argument for category (-c or --category)
    parser.add_argument("-v", "--version", type=int,
                        choices=[4, 5], help="Set EMQX Version")

    args = parser.parse_args()

    if args.version == 4:
        ver = 4

    # create a dashboard
    dashboard = create_dashboard(
        title=title,
        uid=uid,
        templating=create_nodes_template_list(),
    )

    emqx_client_events_dashboard(dashboard, ver)

    dashboard.auto_panel_ids()
    print(dashboard.generate_dashboard())
