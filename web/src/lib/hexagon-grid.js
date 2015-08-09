// Hex math defined here: http://blog.ruslans.com/2011/02/hexagonal-grid-math.html

class HexagonGrid {
  static sign(p1, p2, p3) {
    return (p1.x - p3.x) * (p2.y - p3.y) - (p2.x - p3.x) * (p1.y - p3.y);
  }

  // TODO: Replace with optimized barycentric coordinate method
  static isPointInTriangle(pt, v1, v2, v3) {
    var b1, b2, b3;
    b1 = this.sign(pt, v1, v2) < 0.0;
    b2 = this.sign(pt, v2, v3) < 0.0;
    b3 = this.sign(pt, v3, v1) < 0.0;
    return ((b1 == b2) && (b2 == b3));
  }

  constructor(canvas, options) {
    var radius = options.radius;

    this.radius = radius;

    this.flatTop = false;

    if (this.flatTop) {
      this.height = Math.sqrt(3) * radius;
      this.width = 2 * radius;
    } else {
      this.height = 2 * radius;
      this.width = Math.sqrt(3) * radius;
    }

    this.side = (3 / 2) * radius;
    this.canvas = canvas;
    this.context = this.canvas.getContext('2d');
    this.canvasOriginX = 0;
    this.canvasOriginY = 0;

    this.defaultCell = {
      stroke: "#000",
      fill: "#ddd"
    }
  }

  drawHexGrid(options) {
    let board = options.board || {};

    this.canvasOriginX = options.originX || 0;
    this.canvasOriginY = options.originY || 0;

    for (let col = 0; col < options.columns; col++) {
      for (let row = 0; row < options.rows; row++) {
        let origin = this.getHexOrigin(getIfOffset(col, row));

        let hexPosition = this.getHexPosition(col, row, origin.x, origin.y);

        let cell = board['' + col + row] || this.defaultCell;
        this.drawHex(hexPosition.x, hexPosition.y, cell);
      }
    }
  }

  getHexPosition(column, row, canvasOriginX, canvasOriginY) {
    if (this.flatTop) {
      return {
        x: column * this.side + canvasOriginX,
        y: (row * this.height) + canvasOriginY
      }
    }

    return {
      x: (column * this.width) + canvasOriginX,
      y: row * this.side + canvasOriginY
    }
  }

  getHexOrigin(offset) {
    if (this.flatTop) {
      return {
        x: this.canvasOriginX,
        y: this.canvasOriginY + (offset ? (this.height * 0.5) : 0)
      };
    }

    return {
      x: this.canvasOriginX + (offset ? (this.width * 0.5) : 0),
      y: this.canvasOriginY
    };
  }

  getIfOffset(col, row) {
    if (this.flatTop) {
      return col % 2 !== 0;
    }

    return row % 2 !== 0;
  }

  drawHex(x0, y0, cell) {
    var fillColor = cell.fill;

    this.context.strokeStyle = cell.stroke;
    this.context.beginPath();

    if (this.flatTop) {
      this.context.moveTo(x0 + this.width - this.side, y0);
      this.context.lineTo(x0 + this.side, y0);
      this.context.lineTo(x0 + this.width, y0 + (this.height / 2));
      this.context.lineTo(x0 + this.side, y0 + this.height);
      this.context.lineTo(x0 + this.width - this.side, y0 + this.height);
      this.context.lineTo(x0, y0 + (this.height / 2));

    } else {

      this.context.moveTo(x0 + (this.width / 2), y0);
      this.context.lineTo(x0 + this.width, y0 + this.height - this.side);
      this.context.lineTo(x0 + this.width, y0 + this.side);
      this.context.lineTo(x0 + (this.width / 2), y0 + this.height);
      this.context.lineTo(x0, y0 + this.side);
      this.context.lineTo(x0, y0 + this.height - this.side);
    }


    if (fillColor) {
      this.context.fillStyle = fillColor;
      this.context.fill();
    }

    this.context.closePath();
    this.context.stroke();

    if (cell.text) {
      this.context.fillStyle = "#000";
      this.context.fillText(
        cell.text,
        x0 + (this.width / 2) - (this.width / 4),
        y0 + (this.height - 5)
      );
    }
  }

  getDrawY(column, row) {
    if (column % 2 === 0) {
      return (row * this.height) + this.canvasOriginY;
    }

    return (row * this.height) + this.canvasOriginY + (this.height / 2);
  }

  getDrawX(column) {
    return (column * this.side) + this.canvasOriginX;
  }

  getHexagon(e) {
    var mouse = this.getXYfromEvent(e);

    var tile = this.getSelectedTile(mouse.localX, mouse.localY);

    return {
      column: tile.column,
      row: tile.row
    };
  }

  // Recusivly step up to the body to calculate canvas offset.
  getRelativeCanvasOffset() {
    let x = 0;
    let y = 0;
    let layoutElement = this.canvas;

    if (layoutElement.offsetParent) {
      do {
        x += layoutElement.offsetLeft - layoutElement.scrollLeft;
        y += layoutElement.offsetTop - layoutElement.scrollTop;
        layoutElement = layoutElement.offsetParent;
      } while (layoutElement);

      return {
        x: x,
        y: y
      };
    }
  }

  getRow(column, mouseY) {
    if (column % 2 === 0) {
      return Math.floor((mouseY) / this.height);
    }

    return Math.floor(((mouseY + (this.height * 0.5)) / this.height)) - 1;
  }

  getSelectedTile(mouseX, mouseY) {
    var column = Math.floor((mouseX) / this.side);

    var row = this.getRow(column, mouseY);

    // Test if on left side of frame
    if (mouseX > (column * this.side) && mouseX < (column * this.side) + this.width - this.side) {
      // Now test which of the two triangles we are in
      // Top left triangle points
      var p1 = {
        x: column * this.side,
        y: column % 2 === 0 ? row * this.height : (row * this.height) + (this.height / 2)
      };

      var p2 = {
        x: p1.x,
        y: p1.y + (this.height / 2)
      };

      var p3 = {
        x: p1.x + this.width - this.side,
        y: p1.y
      };

      var mousePoint = {
        x: mouseX,
        y: mouseY
      };

      if (this.isPointInTriangle(mousePoint, p1, p2, p3)) {
        column--;
        if (column % 2 !== 0) {
          row--;
        }
      }

      // Bottom left triangle points
      var p4 = p2;
      var p5 = {
        x: p4.x,
        y: p4.y + (this.height / 2)
      };

      var p6 = {
        x: p5.x + (this.width - this.side),
        y: p5.y
      };

      if (this.isPointInTriangle(mousePoint, p4, p5, p6)) {
        column--;
        if (column % 2 === 0) {
          row++;
        }
      }
    }
    return {
      row: row,
      column: column
    };
  }

  getXYfromEvent(e) {
    var mouseX = e.pageX;
    var mouseY = e.pageY;

    var offSet = this.getRelativeCanvasOffset();

    mouseX -= offSet.x;
    mouseY -= offSet.y;

    return {
      localX: mouseX,
      localY: mouseY
    };
  }
}

export default HexagonGrid;
