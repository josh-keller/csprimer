defmodule ParenthesisMatch do
  @moduledoc """
  Documentation for `ParenthesisMatch`.
  """

  @doc """

  """
  def match(str, stack \\ [])
  def match("", []), do: true
  def match("", [_ | _]), do: false
  def match(<<ch>> <> tail, stack) when ch in '<[({' do
    match(tail, [ch | stack])
  end
  def match(<<ch>> <> _, []) when ch in '>])}', do: false
  def match(")" <> tail, [head | stack]) do
    if head == ?( do
      match(tail, stack)
    else
      false
    end
  end

  def match(<<ch>> <> tail, [head | stack]) when ch in '>}]' do
    if ch - head == 2 do
      match(tail, stack)
    else
      false
    end
  end
end
