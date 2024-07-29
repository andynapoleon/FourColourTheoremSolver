import React from "react";
import { NextReactP5Wrapper } from "@p5-wrapper/next";
import sketch, {
  handleColorMap,
  handleResetMap,
  handleDownloadMap,
} from "../sketch";

export default function Canvas() {
  return <NextReactP5Wrapper sketch={sketch} />;
}

export { handleColorMap, handleResetMap, handleDownloadMap };
