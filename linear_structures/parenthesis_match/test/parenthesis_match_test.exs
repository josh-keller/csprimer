defmodule ParenthesisMatchTest do
  use ExUnit.Case
  doctest ParenthesisMatch

  test "match parens" do
    assert ParenthesisMatch.match("<>{}()[]") == true
    assert ParenthesisMatch.match("<{}()[]") == false
    assert ParenthesisMatch.match("<{}()[]>") == true
    assert ParenthesisMatch.match("<{([])}>") == true
    assert ParenthesisMatch.match("<{([])}>)") == false
  end

  test "file path" do
    assert ParenthesisMatch.path("/etc/bar/../foo/baz.txt") == "/etc/foo/baz.txt"
    assert ParenthesisMatch.path("/etc/bar///foo/baz.txt") == "/etc/bar/foo/baz.txt"
  end
end
