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
)

key = subcharts["authorization"]
title = subchart_metrics[key]['title']
uid = subchart_metrics[key]['uid']
panels = OrderedDict(subchart_metrics[key]['panels'])


def emqx_authz_dashboard(dashboard):
    dashboard.add_row(create_row(
        title=title,
        collapsed=False))

    create_panel(
        dashboard,
        create_table(
            title=panels['subchart_authorize_count']['title'],
            span=6,
            transformations=generate_transformations(
                panels['subchart_authorize_count']['targets']),
            overrides=generate_table_overrides(panels['subchart_authorize_count']['targets'])),
        panels['subchart_authorize_count']['targets'],
        instant=True,
        format=panels['subchart_authorize_count']['format'])

    create_panel(
        dashboard,
        create_timeseries(
            title=panels['subchart_authorize_current_exec_rate']['title'],
            span=4),
        panels['subchart_authorize_current_exec_rate']['targets'])

    create_panel(
        dashboard,
        create_timeseries(
            title=panels['subchart_authorize_last_5m_exec_rate']['title'],
            span=4),
        panels['subchart_authorize_last_5m_exec_rate']['targets'])


if __name__ == '__main__':
    # create a authn dashboard
    dashboard = create_dashboard(
        title=title,
        uid=uid,
        templating=create_auth_template_list(),
    )

    emqx_authz_dashboard(dashboard)

    dashboard.auto_panel_ids()
    print(dashboard.generate_dashboard())
