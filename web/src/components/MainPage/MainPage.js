import React, { PropTypes } from 'react';
import ReactDom from 'react-dom';
import document from 'global/document';
import extend from 'xtend';
import HexagonGrid from '../../lib/hexagon-grid';
import styles from './MainPage.css';
import withStyles from '../../decorators/withStyles';
import window from 'global/window';
import xhr from 'xhr';

@withStyles(styles)
class MainPage extends React.Component {

  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired
  };

  state = {frameIndex: 0};

  componentDidMount() {
    this.redrawBoard();

    setInterval(() => {
      if (this.playNext()) {
        this.nextFrame();
      }
    }, 200);

    window.addEventListener('resize', this.redrawBoard.bind(this));
    window.addEventListener('load', this.redrawBoard.bind(this));
  }

  drawBoard(frame) {
    let grid = ReactDom.findDOMNode(this.refs.gameGrid);
    let hexagonGrid = new HexagonGrid(grid, {
      radius: Math.min(
        1200 / frame.Width / 2,
        700 / frame.Height / Math.sqrt(3)
      )
    });

    hexagonGrid.drawHexGrid(
      this.generateBoardOptions(frame)
    );
  }

  generateBoardOptions(data) {
    return {
      columns: data.Width,
      rows: data.Height,
      board: data.Cells.reduce((board, row) => {
        row.forEach(cell => {
          let gridCell = this.props.emptyCell;
          if (cell.Filled) {
            gridCell = this.props.filledCell;
          }

          board['' + cell.X + cell.Y] = extend(gridCell);
        });
        return board;
      }, {})
    };
  }

  nextFrame() {
    var frame = Math.min(
      this.state.frameIndex + 1,
      this.props.gameData.length - 1
    );
    this.setState({
      frameIndex: frame
    });
    this.redrawBoard(frame);
  }

  playNext() {
    if (this.state.paused) {
      return false;
    }

    return this.state.frameIndex < this.props.gameData.length;
  }

  redrawBoard(frameIndex) {
    this.drawBoard(this.props.gameData[frameIndex || this.state.frameIndex]);
  }

  render() {
    this.context.onSetTitle('ICFP!');
    var gridStyle = {
      position: 'relative'
    };

    return (
      <div className="MainPage">
        <div className="MainPage-container">
          <div className="hexagon-container">
            <canvas className="hexagon-game-grid" ref="gameGrid"
                    width={this.props.gridWidth} height={this.props.gridHeight}>
            </canvas>
          </div>
          <div className="panel-container">
            <h4>Game Controls</h4>
            <br/>
            <h5>Current Frame: {this.state.frameIndex}</h5>
            <br/>
            <button onClick={this.nextFrame.bind(this)}>Next Frame</button>
          </div>
        </div>
      </div>
    );
  }
}

MainPage.defaultProps = {
  emptyCell: {
    fill: '#ddd',
    stroke: '#000'
  },
  filledCell: {
    fill: '#f50',
    stroke: '#000'
  },
  gridWidth: 1000,
  gridHeight: 700
};

export default MainPage;
