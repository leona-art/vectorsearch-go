import { createSignal } from "solid-js";
import "./Counter.css";
import { css } from "~/styled-system/css"

export default function Counter() {
  const [count, setCount] = createSignal(0);
  return (
    <button class={css({
      padding: "10px",
      backgroundColor: "blue",
      color: "white",
      border: "none",
      borderRadius: "5px",
      cursor: "pointer",
    })} onClick={() => setCount(count() + 1)} type="button">
      Clicks: {count()}
    </button>
  );
}
