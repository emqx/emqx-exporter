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
    create_timeseries,
)

key = subcharts["rule-engine-rate"]
title = subchart_metrics[key]['title']
uid = subchart_metrics[key]['uid']
panels = OrderedDict(subchart_metrics[key]['panels'])


def emqx_rule_engine_dashboard(dashboard):
    dashboard.add_row(create_row(
        title=title,
        collapsed=False))

    for _panel_key, panel_value in panels.items():
        create_panel(
            dashboard,
            create_timeseries(
                title=panel_value['title'],
                span=6),
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
