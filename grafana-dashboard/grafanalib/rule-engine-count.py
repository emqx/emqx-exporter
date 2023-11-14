#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from collections import OrderedDict
from metrics import subcharts
from submetrics import subchart_metrics
from common import (
    create_auth_template_list,
    create_dashboard,
    create_panel,
    create_row,
    create_table,
    create_timeseries,
    generate_transformations,
    generate_table_overrides,
    thresholds_2_steps,
)

# key = "rule-engine-count"
key = subcharts["rule-engine-count"]
title = subchart_metrics[key]['title']
uid = subchart_metrics[key]['uid']
panels = OrderedDict(subchart_metrics[key]['panels'])


def emqx_rule_engine_dashboard(dashboard):
    dashboard.add_row(create_row(
        title=title,
        collapsed=False))

    for panel_key, panel_value in panels.items():
        if panel_key == 'subchart_rule_engine_execute_count':
            create_panel(
                dashboard,
                create_table(
                    title=panels['subchart_rule_engine_execute_count']['title'],
                    span=12,
                    transformations=generate_transformations(
                        panels['subchart_rule_engine_execute_count']['targets']),
                    overrides=generate_table_overrides(panels['subchart_rule_engine_execute_count']['targets'])),
                panels['subchart_rule_engine_execute_count']['targets'],
                instant=True,
                format=panels['subchart_rule_engine_execute_count']['format'])
        else:
            create_panel(
                dashboard,
                create_timeseries(
                    title=panel_value['title'],
                    span=3,
                    thresholds=thresholds_2_steps),
                panel_value['targets'])


if __name__ == '__main__':
    # create a authn dashboard
    dashboard = create_dashboard(
        title=title,
        uid=uid,
        templating=create_auth_template_list(),
    )

    emqx_rule_engine_dashboard(dashboard)

    dashboard.auto_panel_ids()
    print(dashboard.generate_dashboard())
