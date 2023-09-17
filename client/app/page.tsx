import Canvas from "./components/canvas";
import { Button } from "@mui/material";
import { ColorMap, ResetMap } from "./components/canvas";

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
      <Button variant="outlined" onClick={ColorMap} className="mt-5">
        Color
      </Button>
      <Button variant="outlined" onClick={ResetMap} className="mt-5 ml-5">
        Reset
      </Button>
    </div>
  );
}
