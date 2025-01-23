let h = 500;
let w = 500;
let grid_h = 400;
let grid_w = 400;
let grid_margin = 50;

let lines = [];
let start, end;
let dragging = false;
let captureImage = false;
let downloadImage = false;
let captured_image = false;
let matrix;

function sketch(p) {
  p.setup = function () {
    p.createCanvas(w, h);
    p.noSmooth();
    p.frameRate(30);
    start = { x: w / 2, y: h / 2 };
  };

  p.draw = function () {
    if (captured_image === false) {
      p.background(200);
      p.stroke(0);
      p.fill(255);
      p.rect(grid_margin, grid_margin, grid_w - 1, grid_h - 1);

      // draw lines
      for (let i = 0; i < lines.length; i++) {
        let a = lines[i][0];
        let b = lines[i][1];
        drawLine(p, a.x, a.y, b.x, b.y);
      }

      // draw temporary line while dragging
      if (dragging) drawLine(p, start.x, start.y, end.x, end.y);
    }

    if (captureImage) {
      let img = p.get(grid_margin, grid_margin, grid_w, grid_h);
      img.loadPixels();
      let array_pixels = img.pixels;
      getData(array_pixels, img.width, img.height);
      captureImage = false;
    }

    if (captured_image) {
      displayColoredMap(p, matrix);
    }

    if (downloadImage) {
      saveCanvasAsImage(p);
    }
  };

  p.mousePressed = function () {
    start = { x: p.mouseX, y: p.mouseY };
  };

  p.mouseDragged = function () {
    end = { x: p.mouseX, y: p.mouseY };
    dragging = true;
  };

  p.mouseReleased = function () {
    dragging = false;
    end = { x: p.mouseX, y: p.mouseY };
    lines.push([start, end]);
    start = { x: end.x, y: end.y };
  };
}

function drawPoint(p, x, y) {
  for (let i = 0; i < 3; i++) {
    for (let j = 0; j < 3; j++) {
      let px = x - 1 + i;
      let py = y - 1 + j;
      if (
        px >= grid_margin &&
        px < grid_w + grid_margin &&
        py >= grid_margin &&
        py < grid_h + grid_margin
      ) {
        p.point(px, py);
      }
    }
  }
}

function drawLine(p, x0, y0, x1, y1) {
  let dx = Math.abs(x1 - x0);
  let dy = Math.abs(y1 - y0);
  let sx = x0 < x1 ? 1 : -1;
  let sy = y0 < y1 ? 1 : -1;
  let err = dx - dy;

  while (true) {
    drawPoint(p, x0, y0);
    if (x0 === x1 && y0 === y1) break;
    let e2 = 2 * err;
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

async function getData(array_pixels, w, h) {
  const apiHost = process.env.NEXT_PUBLIC_API_GATEWAY_URL;

  if (!apiHost) {
    throw new Error("API host is not defined in the environment variables");
  }

  console.log("Auth token: ", localStorage.getItem("token"));

  const res = await fetch(`${apiHost}/api/v1/maps/color`, {
    method: "POST",
    body: JSON.stringify({ image: array_pixels, height: h, width: w }),
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
      "Content-Type": "application/json",
    },
  });

  if (!res.ok) {
    console.log(res);
    throw new Error("Failed to fetch data");
  }

  const data = await res.json();
  matrix = data;
  captured_image = true;
}

function displayColoredMap(p, matrix) {
  let img = p.createImage(matrix[0].length, matrix.length);
  img.loadPixels();
  for (let i = 0; i < img.height; i++) {
    for (let j = 0; j < img.width; j++) {
      img.set(
        j,
        i,
        p.color(
          matrix[i][j][0] * 256,
          matrix[i][j][1] * 256,
          matrix[i][j][2] * 256
        )
      );
    }
  }
  img.updatePixels();
  p.image(img, grid_margin, grid_margin);
}

function saveCanvasAsImage(p) {
  let img = p.get(grid_margin, grid_margin, grid_w, grid_h);
  p.save(img, "map_image.png");
  downloadImage = false;
}

export function handleColorMap() {
  captureImage = true;
  console.log("Coloring map");
}

export function handleResetMap() {
  lines = [];
  captured_image = false;
  console.log("Resetting map");
}

export function handleDownloadMap() {
  downloadImage = true;
  console.log("Downloading map image");
}

export default sketch;
