from collections import OrderedDict
from metrics import subcharts

subchart_metrics = {
    # Messages Subchart
    subcharts["messages"]: {
        "title": "Messages Subchart",
        "uid": "messages",
        "panels": {
            "subchart_messages_sent_rate": {
                "title": "Messages Sent Rate",
                "targets": [
                    {
                        "legendFormat": "{{ instance }}",
                        "expr": "sum by(instance) (irate(emqx_messages_sent{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))"
                    }
                ],
                "format": "timeseries"
            },
            "subchart_messages_received_rate": {
                "title": "Messages Received Rate",
                "targets": [
                    {
                        "legendFormat": "{{ instance }}",
                        "expr": "sum by(instance) (irate(emqx_messages_received{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))"
                    }
                ],
                "format": "timeseries"
            },
            "subchart_messages_dropped_rate": {
                "title": "Messages Dropped Rate",
                "targets": [
                    {
                        "legendFormat": "{{ instance }}",
                        "expr": "sum by(instance) (irate(emqx_messages_dropped{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))"
                    }
                ],
                "format": "timeseries"
            },
        }
    },

    # Connect Auth Subchart
    subcharts["authentication"]: {
        "title": "Connect Auth Subchart",
        "uid": "authentication",
        "panels": {
            "subchart_authenticate_count": {
                "title": "Authenticate Count",
                "targets": [
                    {
                        # "legendFormat": "AuthN Max Rate",
                        "legendFormat": "Max Rate",
                        "expr": "max by(node, resource) (emqx_authentication_exec_max_rate{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",


                    },
                    {
                        # "legendFormat": "AuthN Total",
                        "legendFormat": "Total",
                        "expr": "sum by(node, resource) (emqx_authentication_total{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    },
                    {
                        # "legendFormat": "AuthN Allow",
                        "legendFormat": "Allow",
                        "expr": "sum by(node, resource) (emqx_authentication_allow_count{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    },
                    {
                        # "legendFormat": "AuthN Deny last 15m",
                        "legendFormat": "Deny last 15m",
                        "expr": "sum by(node, resource) (irate(emqx_authentication_deny_count{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"}[15m]))",
                        "thresholds": {
                            "mode": "absolute",
                            "steps": [
                                {
                                    "color": "green",
                                    "value": None
                                },
                                {
                                    "color": "red",
                                    "value": 1
                                }
                            ]
                        },
                    },
                ],
                "format": "table",
            },
            "subchart_authenticate_current_exec_rate": {
                "title": "Authenticate Current Exec Rate",
                "targets": [
                    {
                        "legendFormat": "{{node}} {{resource}}",
                        "expr": "sum by(node, resource) (emqx_authentication_exec_rate{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_authenticate_last_5m_exec_rate": {
                "title": "Authenticate Last 5m Exec Rate",
                "targets": [
                    {
                        "legendFormat": "{{node}} {{resource}}",
                        "expr": "sum by(node, resource) (emqx_authentication_exec_last5m_rate{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    }
                ],
                "format": "timeseries"
            },
        }
    },

    # ACL Auth Subchart
    subcharts["authorization"]: {
        "title": "ACL Auth Subchart",
        "uid": "authorization",
        "panels": {
            "subchart_authorize_count": {
                "title": "Authorize Count",
                "targets": [
                    {
                        # "legendFormat": "AuthZ Max Rate",
                        "legendFormat": "Max Rate",
                        "expr": "max by(node, resource) (emqx_authorization_exec_max_rate{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    },
                    {
                        # "legendFormat": "AuthZ Total",
                        "legendFormat": "Total",
                        "expr": "sum by(node, resource) (emqx_authorization_total{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",

                    },
                    {
                        # "legendFormat": "AuthZ Allow",
                        "legendFormat": "Allow",
                        "expr": "sum by(node, resource) (emqx_authorization_allow_count{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",

                    },
                    {
                        # "legendFormat": "AuthZ Deny last 15m",
                        "legendFormat": "Deny last 15m",
                        "expr": "sum by(node, resource) (irate(emqx_authorization_deny_count{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"}[15m]))",

                        "thresholds": {
                            "mode": "absolute",
                            "steps": [
                                {
                                    "color": "green",
                                    "value": None
                                },
                                {
                                    "color": "red",
                                    "value": 1
                                }
                            ]
                        },
                    }
                ],
                "format": "table",
            },
            "subchart_authorize_current_exec_rate": {
                "title": "Authorize Current Exec Rate",
                "targets": [
                    {
                        "legendFormat": "{{node}} {{resource}}",
                        "expr": "sum by(node, resource) (emqx_authorization_exec_rate{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_authorize_last_5m_exec_rate": {
                "title": "Authorize Last 5m Exec Rate",
                "targets": [
                    {
                        "legendFormat": "{{node}} {{resource}}",
                        "expr": "sum by(node, resource) (emqx_authorization_exec_last5m_rate{cluster=\"$cluster\", node=~\"$node\", resource=~\"$resource\"})",
                    }
                ],
                "format": "timeseries"
            },
        }
    },

    # Client Events Subchart
    subcharts["client-events"]: {
        "title": "Client Events Subchart",
        "uid": "client-events",
        "panels": {
            "subchart_client_connection_events": {
                "title": "Client Connection Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} Connection",
                        "expr": "sum by(instance) (irate(emqx_client_connected{cluster=\"$cluster\",instance=~\"$node\"}[$__rate_interval]))",
                    },
                ],
                "format": "timeseries"
            },
            "subchart_client_disconnection_events": {
                "title": "Client Disconnection Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} Disconnection",
                        "expr": "sum by(instance) (irate(emqx_client_disconnected{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_client_sub_events": {
                "title": "Client Sub Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} Subscribe",
                        "expr": "sum by(instance) (irate(emqx_client_subscribe{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    },
                ],
                "format": "timeseries"
            },
            "subchart_client_unsub_events": {
                "title": "Client Unsub Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} Unsubscribe",
                        "expr": "sum by(instance) (irate(emqx_client_unsubscribe{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_client_connect_auth_events": {
                "title": "Client Connect Auth Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} Connection Auth",
                        "expr": "sum by(instance) (irate(emqx_client_authenticate{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    },
                ],
                "format": "timeseries"
            },
            "subchart_client_acl_auth_events_v4": {
                "title": "Client ACL Auth Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} ACL Auth",
                        "expr": "sum by(instance) (irate(emqx_client_check_acl{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_client_acl_auth_events_v5": {
                "title": "Client ACL Auth Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} ACL Auth",
                        "expr": "sum by(instance) (irate(emqx_client_authorize{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_anonymous_auth_events": {
                "title": "Client Anonymous Auth Events",
                "targets": [
                    {
                        "legendFormat": "{{ instance }} Anonymous Auth",
                        "expr": "sum by(instance) (irate(emqx_client_auth_anonymous{cluster=\"$cluster\", instance=~\"$node\"}[$__rate_interval]))",
                    },
                ],
                "format": "timeseries"
            },
        }
    },

    # Rule Engine Rate Subchart
    subcharts["rule-engine-rate"]: {
        "title": "Rule Engine Rate Subchart",
        "uid": "rule-engine-rate",
        "panels": {
            "subchart_rule_engine_current_exec_rate": {
                "title": "Rule Engine Current Exec Rate",
                "targets": [
                    {
                        "legendFormat": "{{ node }} {{ rule }}",
                        "expr": "sum by(node, rule) (emqx_rule_exec_rate{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_rule_engine_last_5m_exec_rate": {
                "title": "Rule Engine Last 5m Exec Rate",
                "targets": [
                    {
                        "legendFormat": "{{ node }}:{{ rule }}",
                        "expr": "sum by(node, rule) (emqx_rule_exec_last5m_rate{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    }
                ],
                "format": "timeseries"
            },
        },
    },

    # Rule Engine Count Subchart
    subcharts["rule-engine-count"]: {
        "title": "Rule Engine Count Subchart",
        "uid": "rule-engine-count",
        "panels": {
            "subchart_rule_engine_execute_count": {
                "title": "Rule Engine Execute Count",
                "targets": [
                    {
                        "legendFormat": "Topic Hit Cout",
                        "expr": "sum by(node, rule) (emqx_rule_topic_hit_count{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    },
                    {
                        "legendFormat": "Max Rate",
                        "expr": "max by(node, rule) (emqx_rule_exec_max_rate{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    },
                    {
                        "legendFormat": "Exec Pass",
                        "expr": "sum by(node, rule) (emqx_rule_exec_pass_count{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    },
                    {
                        "legendFormat": "Exec Exception last 15m",
                        "expr": "sum by(node, rule) (irate(emqx_rule_exec_exception_count{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"}[15m]))",
                        "thresholds": {
                            "mode": "absolute",
                            "steps": [
                                {
                                    "color": "green",
                                    "value": None
                                },
                                {
                                    "color": "red",
                                    "value": 1
                                }
                            ]
                        }
                    },
                    {
                        "legendFormat": "Exec No Result 15m",
                        "expr": "sum by(node, rule) (emqx_rule_exec_no_result_count{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                        "thresholds": {
                            "mode": "absolute",
                            "steps": [
                                {
                                    "color": "green",
                                    "value": None
                                },
                                {
                                    "color": "orange",
                                    "value": 1
                                }
                            ]
                        },
                    },
                    {
                        "legendFormat": "Call Action Total",
                        "expr": "sum by(node, rule) (emqx_rule_action_total{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    },
                    {
                        "legendFormat": "Action Success",
                        "expr": "sum by(node, rule) (emqx_rule_action_success{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"})",
                    },
                    {
                        "legendFormat": "Action Failed last 15m",
                        "expr": "sum by(node, rule) (increase(emqx_rule_action_failed{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"}[15m]))",
                        "thresholds": {
                            "mode": "absolute",
                            "steps": [
                                {
                                    "color": "green",
                                    "value": None
                                },
                                {
                                    "color": "red",
                                    "value": 1
                                }
                            ]
                        },
                    },
                ],
                "format": "table",
            },
            "subchart_rule_engine_exec_success": {
                "title": "Rule Engine Exec Success",
                "targets": [
                    {
                        "legendFormat": "{{ node }}:{{ rule }}",
                        "expr": "sum by(rule) (irate(emqx_rule_exec_pass_count{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_rule_engine_exec_failure": {
                "title": "Rule Engine Exec Exception",
                "targets": [
                    {
                        "legendFormat": "{{ node }}:{{ rule }}",
                        "expr": "sum by(rule) (irate(emqx_rule_exec_exception_count{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"}[$__rate_interval]))",

                    }
                ],
                "format": "timeseries"
            },
            "subchart_rule_engine_action_success": {
                "title": "Rule Engine Action Success",
                "targets": [
                    {
                        "legendFormat": "{{ node }}:{{ rule }}",
                        "expr": "sum by(rule) (irate(emqx_rule_action_success{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
            "subchart_rule_engine_action_failure": {
                "title": "Rule Engine Action Failure",
                "targets": [
                    {
                        "legendFormat": "{{ node }}:{{ rule }}",
                        "expr": "sum by(rule) (irate(emqx_rule_action_success{cluster=\"$cluster\", node=~\"$node\", rule=~\"$rule\"}[$__rate_interval]))",
                    }
                ],
                "format": "timeseries"
            },
        }
    },
}
