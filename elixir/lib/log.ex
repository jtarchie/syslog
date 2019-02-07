defmodule SyslogLog do
  @type t :: %__MODULE__{
          version: integer(),
          severity: integer(),
          facility: integer(),
          timestamp: DateTime.t() | nil,
          hostname: String.t() | nil,
          app_name: String.t() | nil,
          proc_id: String.t() | nil,
          msg_id: String.t() | nil,
          structure_data: any(),
          message: String.t() | nil
        }

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
