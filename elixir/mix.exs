defmodule Syslog.MixProject do
  use Mix.Project

  def project do
    [
      app: :syslog,
      version: "0.1.0",
      elixir: "~> 1.9-dev",
      start_permanent: Mix.env() == :prod,
      deps: deps()
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
      {:mix_test_watch, "~> 0.8", only: :dev, runtime: false}
    ]
  end
end
