defmodule Syslog.MixProject do
  use Mix.Project

  def project do
    [
      app: :syslog,
      version: "0.1.0",
      elixir: "~> 1.8",
      start_permanent: Mix.env() == :prod,
      deps: deps(),
      dialyzer: [flags: [:unmatched_returns,:error_handling,:race_conditions, :no_opaque]]
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  def deps do
    [
      {:nimble_parsec, "~> 0.2"},
      {:mix_test_watch, "~> 0.8", only: :dev, runtime: false},
      {:benchee, "~> 0.13", only: :dev},
      {:exprof, "~> 0.2.0", only: :dev},
      {:dialyxir, "~> 1.0.0-rc.4", only: [:dev], runtime: false}
    ]
  end
end
