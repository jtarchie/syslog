defmodule Syslog do
  @compile [:native, {:hipe, [:verbose, :o3]}]
  use Bitwise

  defmodule Element do
    @type t :: %__MODULE__{
            id: String.t(),
            properties: list(Property.t())
          }
    defstruct [:id, :properties]
  end

  defmodule Property do
    @type t :: %__MODULE__{
            key: String.t(),
            value: String.t()
          }
    defstruct [:key, :value]
  end

  @spec parse(String.t()) :: {:ok, SyslogLog.t()} | {:error, String.t()}
  def parse(msg) do
    case SyslogParser.message(msg) do
      {:ok, val, _, _, _, _} ->
        SyslogParser.message(msg)
        log = build(%SyslogLog{structure_data: []}, val)
        {:ok, log}

      {:error, msg, _, _, _, _} ->
        {:error, msg}
    end
  end

  defp build(log, [{:version, version} | p]) do
    log = %{log | version: :erlang.list_to_integer(version)}
    build(log, p)
  end

  defp build(log, [{:prival, [prival]} | p]) do
    log = %{log | severity: prival &&& 7, facility: prival >>> 3, priority: prival}
    build(log, p)
  end

  defp build(log, [{:datetime, [timestamp_string]} | p]) do
    timestamp =
      case DateTime.from_iso8601(timestamp_string) do
        {:ok, timestamp, _} ->
          timestamp

        _ ->
          nil
      end

    log = %{log | timestamp: timestamp}
    build(log, p)
  end

  defp build(log, [{:hostname, [hostname]} | p]) do
    log = %{log | hostname: hostname}
    build(log, p)
  end

  defp build(log, [{:app_name, [app_name]} | p]) do
    log = %{log | app_name: app_name}
    build(log, p)
  end

  defp build(log, [{:proc_id, [proc_id]} | p]) do
    log = %{log | proc_id: proc_id}
    build(log, p)
  end

  defp build(log, [{:msg_id, [msg_id]} | p]) do
    log = %{log | msg_id: msg_id}
    build(log, p)
  end

  defp build(log, [{:message, [message]} | p]) do
    log = %{log | message: message}
    build(log, p)
  end

  defp build(log, [{:sd_element, sd_element} | p]) do
    element = build_sd_element(%Element{properties: []}, sd_element)
    log = %{log | structure_data: [element | log.structure_data]}
    build(log, p)
  end

  defp build(log, [_ | p]) do
    build(log, p)
  end

  defp build(log, []) do
    log
  end

  defp build_sd_element(element, [{:sd_id, [id]} | properties]) do
    element = %{element | id: id}
    build_sd_element(element, properties)
  end

  defp build_sd_element(element, [{:sd_param, [key, value]} | properties]) do
    element = %{element | properties: [%Property{key: key, value: value} | element.properties]}
    build_sd_element(element, properties)
  end

  defp build_sd_element(element, []) do
    element
  end
end
