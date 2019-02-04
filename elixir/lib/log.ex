defmodule SyslogLog do
  defstruct [
    :version,
    :severity,
    :facility,
    :priority,
    :timestamp,
    :hostname,
    :app_name,
    :proc_id,
    :msg_id,
    :structure_data,
    :message
  ]
end
