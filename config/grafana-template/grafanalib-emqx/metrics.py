#!/usr/bin/env python3
# -*- coding: utf-8 -*-

metrics = {
    # General
    "general_row": {
        "title": "General",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 0
        },
    },
    "cluster_status": {
        "title": "Cluster Status",
        "targets": [
            {
                "legendFormat": "Status",
                "expr": "emqx_cluster_status{cluster=\"$cluster\"}"
            }
        ],
        "format": "stat",
    },
    "license": {
        "title": "License",
        "targets": [
            {
                "legendFormat": "Days Remaining",
                "expr": "sum(emqx_license_remaining_days{cluster=\"$cluster\"})",
                "thresholds": {
                    "mode": "absolute",
                    "steps": [
                        {
                            "color": "red",
                            "value": None
                        },
                        {
                            "color": "#EAB839",
                            "value": 30
                        },
                        {
                            "color": "green",
                            "value": 90
                        }
                    ]
                }
            },
            {
                "legendFormat": "Expiry At",
                "expr": "sum(emqx_license_expiration_time{cluster=\"$cluster\"})",
                "datetime": True
            },
            {
                "legendFormat": "Max Conns",
                "expr": "sum(emqx_license_max_client_limit{cluster=\"$cluster\"})"
            }
        ],
        "format": "table",
        "gridPos": {
            "h": 3,
            "w": 6,
            "x": 0,
            "y": 4
        },
    },
    "active_connections": {
        "title": "Active Connections",
        "targets": [
            {
                "legendFormat": "Connections",
                "expr": "sum(emqx_connections_count{instance=~\".*\", cluster=\"$cluster\"})"
            },
            # {
            #     "legendFormat": "Max Connections",
            #     "expr": "sum(emqx_license_max_client_limit{cluster=\"$cluster\"})",
            # },
        ],
        "format": "gauge"
        # "format": "timeseries",
    },
    "cluster_message_rate": {
        "title": "Cluster Message Rate",
        "targets": [
            {
                "legendFormat": "Msg Input Period Second",
                "expr": "emqx_messages_input_period_second{cluster=\"$cluster\"}"
            },
            {
                "legendFormat": "Msg Output Period Second",
                "expr": "emqx_messages_output_period_second{cluster=\"$cluster\"}"
            }
        ],
        "format": "timeseries",
        "gridPos": {
            "h": 6,
            "w": 5,
            "x": 10,
            "y": 1
        },
    },
    "nodes_running": {
        "title": "Nodes Running",
        "targets": [
            {
                "legendFormat": "Running",
                "expr": "max(emqx_cluster_nodes_running{instance=~\".*\", cluster=\"$cluster\"})",
                "color": "green",
            },
            {
                "legendFormat": "Stopped",
                "expr": "max(emqx_cluster_nodes_stopped{instance=~\".*\", cluster=\"$cluster\"})",
                "color": "dark-red",
            }
        ],
        "format": "timeseries",
        "gridPos": {
            "h": 6,
            "w": 5,
            "x": 15,
            "y": 1
        },
    },
    "exporter_latency": {
        "title": "Exporter Latency",
        "targets": [
            {
                "legendFormat": "Latency",
                "expr": "sum(emqx_scrape_collector_duration_seconds{cluster=\"$cluster\"})"
            }
        ],
        "format": "timeseries",
        "gridPos": {
            "h": 6,
            "w": 4,
            "x": 20,
            "y": 1
        },
    },
    "sessions": {
        "title": "Sessions",
        "targets": [
            {
                "legendFormat": "Max",
                "expr": "sum(emqx_license_max_client_limit{cluster=\"$cluster\"})"
            },
            {
                "legendFormat": "Count",
                "expr": "sum(emqx_sessions_count{instance=~\".*\", cluster=\"$cluster\"})"
            },
            {
                "legendFormat": "{{ instance }}",
                "expr": "sum by(instance) (emqx_sessions_count{cluster=\"$cluster\"})"
            },
        ],
        "format": "timeseries"
    },
    "connections": {
        "title": "Connections",
        "targets": [
            {
                "legendFormat": "Total",
                "expr": "sum(emqx_live_connections_count{instance=~\".*\", cluster=\"$cluster\"})"
            },
            {
                "legendFormat": "{{ instance }}",
                "expr": "sum by(instance) (emqx_connections_count{cluster=\"$cluster\"})"
            }
        ],
        "format": "timeseries"
    },
    "subscriptions": {
        "title": "Subscriptions",
        "targets": [
            {
                "legendFormat": "Subscriptions",
                "expr": "sum(emqx_suboptions_count{instance=~\".*\", cluster=\"$cluster\"})"
            },
            {
                "legendFormat": "{{ instance }}",
                "expr": "sum by(instance) (emqx_suboptions_count{cluster=\"$cluster\"})"
            }
        ],
        "format": "timeseries"
    },
    "rule_engine_last_5m_exec_rate": {
        "title": "Rule Engine Last 5m Exec Rate",
        "targets": [
            {
                "legendFormat": "{{ rule }}",
                "expr": "sum by(rule) (emqx_rule_exec_last5m_rate{cluster=\"$cluster\", node=~\".*\"})"
            }
        ],
        "format": "timeseries"
    },
    "connect_auth": {
        "title": "Connect Auth",
        "targets": [
            {
                "legendFormat": "Auth Success",
                "expr": "sum(irate(emqx_packets_connect{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval])) - sum(irate(emqx_packets_connack_error {instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "Auth Failure",
                "expr": "sum(irate(emqx_packets_connack_auth_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "acl_auth": {
        "title": "ACL Auth",
        "targets": [
            {
                "legendFormat": "Publish ACL Failure",
                "expr": "sum(irate(emqx_packets_publish_auth_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "Sub ACL Failure ",
                "expr": "sum(irate(emqx_packets_subscribe_auth_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "data_bridge_queuing": {
        "title": "Data Bridge Queuing",
        "targets": [
            {
                "legendFormat": "{{type}}-{{name}}",
                "expr": "sum by(type, name) (emqx_rule_bridge_queuing{cluster=\"$cluster\"})"
            }
        ],
        "format": "timeseries"
    },
    "rule_engine_current_exec_rate": {
        "title": "Rule Engine Current Exec Rate",
        "targets": [
            {
                "legendFormat": "{{ rule }}",
                "expr": "sum by(rule) (emqx_rule_exec_rate{cluster=\"$cluster\", node=~\".*\"})"
            }
        ],
        "format": "timeseries"
    },
    "rule_engine_exec_success": {
        "title": "Rule Engine Exec Success",
        "targets": [
            {
                "legendFormat": "{{rule}}",
                "expr": "sum by(rule) (irate(emqx_rule_exec_pass_count{cluster=\"$cluster\", node=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "rule_engine_exec_failure": {
        "title": "Rule Engine Exec Failure",
        "targets": [
            {
                "legendFormat": "{{rule}}",
                "expr": "sum by(rule) (irate(emqx_rule_exec_failure_count{cluster=\"$cluster\", node=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "rule_engine_action_success": {
        "title": "Rule Engine Action Success",
        "targets": [
            {
                "legendFormat": "{{rule}}",
                "expr": "sum by(rule) (irate(emqx_rule_action_success{cluster=\"$cluster\", node=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "rule_engine_action_failure": {
        "title": "Rule Engine Action Failure",
        "targets": [
            {
                "legendFormat": "{{rule}}",
                "expr": "sum by(rule) (irate(emqx_rule_action_failed{cluster=\"$cluster\", node=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    # Event
    "event_row": {
        "title": "Event",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 29
        },
    },
    "client_connection_events": {
        "title": "Client Connection Events",
        "targets": [
            {
                "legendFormat": "Connections",
                "expr": "sum(irate(emqx_client_connected{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "Disconnections",
                "expr": "sum(irate(emqx_client_disconnected{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "client_sub_events": {
        "title": "Client Sub Events",
        "targets": [
            {
                "legendFormat": "Subscribes",
                "expr": "sum(irate(emqx_client_subscribe{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "Unsubscribes",
                "expr": "sum(irate(emqx_client_unsubscribe{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "client_connect_auth_events": {
        "title": "Client Connect Auth Events",
        "targets": [
            {
                "legendFormat": "Connect Auth",
                "expr": "sum(irate(emqx_client_authenticate{cluster=\"$cluster\",instance=~\".*\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "Anonymous Auth",
                "expr": "sum(irate(emqx_client_auth_anonymous{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "client_acl_auth_events": {
        "title": "Client ACL Auth Events",
        "targets": [
            {
                "legendFormat": "ACL Auth",
                "expr": "sum(irate(emqx_client_authorize{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    # Packets
    "packets_row": {
        "title": "Packets",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 30
        },
    },
    "packets_connections": {
        "title": "Packets Connections",
        "targets": [
            {
                "legendFormat": "Packets Connect",
                "expr": "sum(irate(emqx_packets_connect{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "green",
            },
            {
                "legendFormat": "Packets Connack Sent",
                "expr": "sum(irate(emqx_packets_connack_sent{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "blue",
            },
            {
                "legendFormat": "Packets Connack Error",
                "expr": "sum(irate(emqx_packets_connack_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "red",
            }
        ],
        "format": "timeseries"
    },
    "packets_disconnections": {
        "title": "Packets Disconnections",
        "targets": [
            {
                "legendFormat": "Packets Disconnect Sent",
                "expr": "sum(irate(emqx_packets_disconnect_sent{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "semi-dark-red",
            },
            {
                "legendFormat": "Packets Disconnect Received",
                "expr": "sum(irate(emqx_packets_disconnect_received{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "super-light-red",
            }
        ],
        "format": "timeseries"
    },
    "packets_publish": {
        "title": "Packets Publish",
        "targets": [
            {
                "legendFormat": "Packets Publish Sent",
                "expr": "sum(irate(emqx_packets_publish_sent{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "blue",
            },
            {
                "legendFormat": "Packets Publish Received",
                "expr": "sum(irate(emqx_packets_publish_received{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "dark-purple",
            },
            {
                "legendFormat": "Packets Publish Dropped",
                "expr": "sum(irate(emqx_packets_publish_dropped{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "super-light-red",
            },
            {
                "legendFormat": "Packets Publish Error",
                "expr": "sum(irate(emqx_packets_publish_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "dark-red",
            }
        ],
        "format": "timeseries"
    },
    "packets_subscribe_and_unsubscribe": {
        "title": "Packets Subscribe/Unsubscribe",
        "targets": [
            {
                "legendFormat": "Packets Subscribe Received",
                "expr": "sum(irate(emqx_packets_subscribe_received{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "dark-purple",
            },
            {
                "legendFormat": "Packets Suback Sent",
                "expr": "sum(irate(emqx_packets_suback_sent{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "blue",
            },
            {
                "legendFormat": "Packets Subscribe Error",
                "expr": "sum(irate(emqx_packets_subscribe_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "dark-red",
            },
            {
                "legendFormat": "Packets Unsubscribe Received",
                "expr": "sum(irate(emqx_packets_unsubscribe_received{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "orange",
            },
            {
                "legendFormat": "Packets Unsubscribe Error",
                "expr": "sum(irate(emqx_packets_unsubscribe_error{instance=~\".*\", cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "super-light-red",
            }
        ],
        "format": "timeseries"
    },
    # Messages
    "messages_row": {
        "title": "Messages",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 31
        },
    },
    "messages_count": {
        "title": "Messages Count",
        "targets": [
            {
                "legendFormat": "Messages Received",
                "expr": "sum(irate(emqx_messages_received{instance=~\".*\",cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "dark-purple"
            },
            {
                "legendFormat": "Messages Sent",
                "expr": "sum(irate(emqx_messages_sent{instance=~\".*\",cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "blue"
            },
            {
                "legendFormat": "Messages Dropped",
                "expr": "sum(irate(emqx_messages_dropped{instance=~\".*\",cluster=\"$cluster\"}[$__rate_interval]))",
                "color": "super-light-red"
            }
        ],
        "format": "timeseries"
    },
    "messages_qos_received": {
        "title": "Messages QOS Received",
        "targets": [
            {
                "legendFormat": "QOS0",
                "expr": "sum(irate(emqx_messages_qos0_received{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "QOS1",
                "expr": "sum(irate(emqx_messages_qos1_received{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "QOS2",
                "expr": "sum(irate(emqx_messages_qos2_received{cluster=\"$cluster\", instance=~\".*\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    "cluster_traffic_statistics": {
        "title": "Cluster Traffic Statistics",
        "targets": [
            {
                "legendFormat": "Bytes received",
                "expr": "sum(irate(emqx_bytes_received{instance=~\".*\",cluster=\"$cluster\"}[$__rate_interval]))"
            },
            {
                "legendFormat": "Bytes sent",
                "expr": "sum(irate(emqx_bytes_sent{instance=~\".*\",cluster=\"$cluster\"}[$__rate_interval]))"
            }
        ],
        "format": "timeseries"
    },
    # Rule Engine
    "rule_engine_row": {
        "title": "Rule Engine",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 32
        },
    },
    "data_bridge_status": {
        "title": "Data Bridge Status",
        "targets": [
            {
                "legendFormat": "Status",
                "expr": "sum by(type, name) (emqx_rule_bridge_status{cluster=\"$cluster\"})",
                "mappings": [
                    {
                        "options": {
                            "1": {
                                "color": "red",
                                "index": 1,
                                "text": "Disconected"
                            },
                            "2": {
                                "color": "green",
                                "index": 0,
                                "text": "Connected"
                            }
                        },
                        "type": "value"
                    },
                    {
                        "options": {
                            "match": None,
                            "result": {
                                "color": "red",
                                "index": 2,
                                "text": "N/A"
                            }
                        },
                        "type": "special"
                    }
                ]
            },
            {
                "legendFormat": "Queuing",
                "expr": "sum by(type, name) (emqx_rule_bridge_queuing{cluster=\"$cluster\"})",
                "thresholds": {
                    "mode": "absolute",
                    "steps": [
                        {
                            "color": "green",
                            "value": None,
                        },
                        {
                            "color": "red",
                            "value": 1,
                        }
                    ]
                }
            }
        ],
        "format": "table"
    },
    "rule_engine_execute_count": {
        "title": "Rule Engine Execute Count",
        "targets": [
            {
                "legendFormat": "Topic Hit Cout",
                "expr": "sum by(rule) (emqx_rule_topic_hit_count{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                "legendFormat": "Max Rate",
                "expr": "max by(rule) (emqx_rule_exec_max_rate{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                "legendFormat": "Exec Pass",
                "expr": "sum by(rule) (emqx_rule_exec_pass_count{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                "legendFormat": "Exec Failed last 15m",
                "expr": "sum by(rule) (irate(emqx_rule_exec_failure_count{cluster=\"$cluster\", node=~\".*\"}[15m]))",
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
                "legendFormat": "Exec No Result",
                "expr": "sum by(rule) (emqx_rule_exec_no_result_count{cluster=\"$cluster\", node=~\".*\"})",
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
                "expr": "sum by(rule) (emqx_rule_action_total{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                "legendFormat": "Action Success",
                "expr": "sum by(rule) (emqx_rule_action_success{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                "legendFormat": "Action Failed last 15m",
                "expr": "sum by(rule) (increase(emqx_rule_action_failed{cluster=\"$cluster\", node=~\".*\"}[15m]))",
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
    # Connect Auth
    "connect_auth_row": {
        "title": "Connect Auth",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 33
        },
    },
    "authenticate_count": {
        "title": "Authenticate Count",
        "targets": [
            {
                # "legendFormat": "AuthN Status",
                "legendFormat": "Status",
                "expr": "sum by(resource) (emqx_authentication_resource_status{cluster=\"$cluster\"})",
                "mappings": [
                    {
                        "options": {
                            "1": {
                                "color": "red",
                                "index": 1,
                                "text": "Disconnected"
                            },
                            "2": {
                                "color": "green",
                                "index": 0,
                                "text": "Connected"
                            }
                        },
                        "type": "value"
                    },
                    {
                        "options": {
                            "match": "null",
                            "result": {
                                "color": "red",
                                "index": 2,
                                "text": "Unknown"
                            }
                        },
                        "type": "special"
                    }
                ],
            },
            {
                # "legendFormat": "AuthN Max Rate",
                "legendFormat": "Max Rate",
                "expr": "max by(resource) (emqx_authentication_exec_max_rate{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                # "legendFormat": "AuthN Total",
                "legendFormat": "Total",
                "expr": "sum by(resource) (emqx_authentication_total{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                # "legendFormat": "AuthN Allow",
                "legendFormat": "Allow",
                "expr": "sum by(resource) (emqx_authentication_allow_count{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                # "legendFormat": "AuthN Deny last 15m",
                "legendFormat": "Deny last 15m",
                "expr": "sum by(resource) (irate(emqx_authentication_deny_count{cluster=\"$cluster\", node=~\".*\"}[15m]))",
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
    "authenticate_current_exec_rate": {
        "title": "Authenticate Current Exec Rate",
        "targets": [
            {
                "legendFormat": "{{resource}}",
                "expr": "sum by(resource) (emqx_authentication_exec_rate{cluster=\"$cluster\", node=~\".*\"})"
            }
        ],
        "format": "timeseries"
    },
    "authenticate_last_5m_exec_rate": {
        "title": "Authenticate Last 5m Exec Rate",
        "targets": [
            {
                "legendFormat": "{{resource}}",
                "expr": "sum by(resource) (emqx_authentication_exec_last5m_rate{cluster=\"$cluster\", node=~\".*\"})"
            }
        ],
        "format": "timeseries"
    },
    # ACL Auth
    "acl_auth_row": {
        "title": "ACL Auth",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 34
        },
    },
    "authorize_count": {
        "title": "Authorize Count",
        "targets": [
            {
                # "legendFormat": "AuthZ Status",
                "legendFormat": "Status",
                "expr": "sum by(resource) (emqx_authorization_resource_status{cluster=\"$cluster\"})",
                "mappings": [
                    {
                        "options": {
                            "1": {
                                "color": "red",
                                "index": 1,
                                "text": "Disconnected"
                            },
                            "2": {
                                "color": "green",
                                "index": 0,
                                "text": "Connected"
                            }
                        },
                        "type": "value"
                    },
                    {
                        "options": {
                            "match": None,
                            "result": {
                                "color": "red",
                                "index": 2,
                                "text": "Unknown"
                            }
                        },
                        "type": "special"
                    }
                ]
            },
            {
                # "legendFormat": "AuthZ Max Rate",
                "legendFormat": "Max Rate",
                "expr": "max by(resource) (emqx_authorization_exec_max_rate{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                # "legendFormat": "AuthZ Total",
                "legendFormat": "Total",
                "expr": "sum by(resource) (emqx_authorization_total{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                # "legendFormat": "AuthZ Allow",
                "legendFormat": "Allow",
                "expr": "sum by(resource) (emqx_authorization_allow_count{cluster=\"$cluster\", node=~\".*\"})",
            },
            {
                # "legendFormat": "AuthZ Deny last 15m",
                "legendFormat": "Deny last 15m",
                "expr": "sum by(resource) (irate(emqx_authorization_deny_count{cluster=\"$cluster\", node=~\".*\"}[15m]))",
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
    "authorize_current_exec_rate": {
        "title": "Authorize Current Exec Rate",
        "targets": [
            {
                "legendFormat": "{{resource}}",
                "expr": "sum by(resource) (emqx_authorization_exec_rate{cluster=\"$cluster\", node=~\".*\"})"
            }
        ],
        "format": "timeseries"
    },
    "authorize_last_5m_exec_rate": {
        "title": "Authorize Last 5m Exec Rate",
        "targets": [
            {
                "legendFormat": "{{resource}}",
                "expr": "sum by(resource) (emqx_authorization_exec_last5m_rate{cluster=\"$cluster\", node=~\".*\"})"
            }
        ],
        "format": "timeseries"
    },
    # System
    "system_row": {
        "title": "System",
        "gridPos":  {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 35
        },
    },
    "last_1m_cpu_load": {
        "title": "Last 1m CPU Load",
        "targets": [
            {
                "legendFormat": "{{node}}",
                "expr": "sum by(node) (emqx_cluster_cpu_load{cluster=\"$cluster\", load=\"load1\"})"
            }
        ],
        "format": "timeseries"
    },
    "last_5m_cpu_load": {
        "title": "Last 5m CPU Load",
        "targets": [
            {
                "legendFormat": "{{node}}",
                "expr": "sum by(node) (emqx_cluster_cpu_load{cluster=\"$cluster\", load=\"load5\"})"
            }
        ],
        "format": "timeseries"
    },
    "last_15m_cpu_load": {
        "title": "Last 15m CPU Load",
        "targets": [
            {
                "legendFormat": "{{node}}",
                "expr": "sum by(node) (emqx_cluster_cpu_load{cluster=\"$cluster\", load=\"load15\"})"
            }
        ],
        "format": "timeseries"
    }
}
