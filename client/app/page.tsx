import Canvas from "./components/canvas";
import React from "react";

export default function Home() {
  return (
    <div
      style={{
        position: "absolute",
        left: "50%",
        top: "50%",
        transform: "translate(-50%, -50%)",
      }}
    >
      <Canvas />
    </div>
  );
}
