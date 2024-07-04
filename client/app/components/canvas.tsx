"use client";

import { NextReactP5Wrapper } from "@p5-wrapper/next";
import { P5WrapperClassName, type Sketch } from "@p5-wrapper/react";
import { Life_Savers } from "next/font/google";

var h = 300;
var w = 500;
var grid_h = 200;
var grid_w = 400;
var grid_margin = 50;

var lines: { x: any; y: any }[][] = [];
var start: { x: any; y: any }, end: { x: any; y: any }; // JSON
var dragging: boolean; // boolean
var captured_image = false;
var captureImage = false;
var matrix: any;

const sketch: Sketch = (p5) => {
  p5.setup = () => {
    p5.createCanvas(w, h);
    p5.noSmooth();
    p5.frameRate(30);
    dragging = false;
    start = {
      x: w / 2,
      y: h / 2,
    };
    p5.loop();
  };

  p5.draw = () => {
    if (captured_image === false) {
      p5.background(200);
      p5.stroke(0);
      p5.fill(255);
      p5.rect(50, 50, grid_w - 1, grid_h - 1);

      // draw lines
      let i = 0;
      for (i = 0; i < lines.length; i++) {
        var a = lines[i][0];
        var b = lines[i][1];
        draw_line(a.x, a.y, b.x, b.y);
      }

      // draw temporary line while dragging
      if (dragging) draw_line(start.x, start.y, end.x, end.y);
    }

    // save image
    if (captureImage) {
      let img = p5.get(50, 50, grid_w, grid_h);
      //p5.save(img, "captured_image.png");
      captureImage = false; // reset captureImage to false
      img.loadPixels();
      var array_pixels = img.pixels;
      getData(array_pixels, img.width, img.height);
      console.log(matrix);
    }

    if (captured_image) {
      displayColoredMap(matrix);
    }
  };

  async function displayColoredMap(matrix: any) {
    // var height = matrix[0].length;
    // var width = matrix.length;
    let img = p5.createImage(matrix[0].length, matrix.length);
    img.loadPixels();
    for (let i = 0; i < img.height; i++) {
      for (let j = 0; j < img.width; j++) {
        img.set(
          j,
          i,
          p5.color(
            matrix[i][j][0] * 256,
            matrix[i][j][1] * 256,
            matrix[i][j][2] * 256
          )
        );
      }
    }
    console.log(img.pixels);
    img.updatePixels();
    img.loadPixels();
    p5.image(img, 50, 50);
  }

  async function getData(array_pixels: any, w: any, h: any) {
    const apiHost = process.env.NEXT_PUBLIC_API_URL;

    if (!apiHost) {
      throw new Error("API host is not defined in the environment variables");
    }

    const res = await fetch(`${apiHost}/api/solve`, {
      method: "POST",
      body: JSON.stringify({ image: array_pixels, height: h, width: w }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!res.ok) {
      // This will activate the closest `error.js` Error Boundary
      throw new Error("Failed to fetch data");
    }

    const data = await res.json();

    console.log(data);
    matrix = data;
    captured_image = true;
  }

  function draw_point(x: number, y: number) {
    var i = 0;
    var j = 0;
    for (i = 0; i < 3; i++) {
      for (j = 0; j < 3; j++) {
        var px = x - 1 + i;
        var py = y - 1 + j;
        if (
          px >= grid_margin &&
          px < grid_w + grid_margin &&
          py >= grid_margin &&
          py < grid_h + grid_margin
        ) {
          p5.point(px, py);
        }
      }
    }
  }

  function draw_line(x0: number, y0: number, x1: number, y1: number) {
    var dx = Math.abs(x1 - x0);
    var dy = Math.abs(y1 - y0);
    var sx = x0 < x1 ? 1 : -1;
    var sy = y0 < y1 ? 1 : -1;
    var err = dx - dy;

    while (true) {
      draw_point(x0, y0);
      if (x0 == x1 && y0 == y1) break;
      var e2 = 2 * err;
      if (e2 > -dy) {
        err -= dy;
        x0 += sx;
      }
      if (e2 < dx) {
        err += dx;
        y0 += sy;
      }
    }
  }

  p5.mousePressed = () => {
    start = {
      x: p5.mouseX,
      y: p5.mouseY,
    };
    //console.log(start);
  };

  p5.mouseDragged = () => {
    end = {
      x: p5.mouseX,
      y: p5.mouseY,
    };
    //console.log(end);
    dragging = true;
    //console.log(dragging);
  };

  p5.mouseReleased = () => {
    dragging = false;
    end = {
      x: p5.mouseX,
      y: p5.mouseY,
    };
    lines.push([start, end]);
    start = {
      x: end.x,
      y: end.y,
    };
  };
};

export function ColorMap() {
  captureImage = true;
  console.log("Hello Color");
}

export function ResetMap() {
  lines = [];
  captured_image = false;
  console.log("Hello Reset");
}

export default function Canvas() {
  return <NextReactP5Wrapper sketch={sketch} />;
}
