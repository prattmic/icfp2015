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
    this.setState({
      mounted: true
    });

    this.redrawBoard();

    setInterval(() => {
      if (this.playNext()) {
        this.nextFrame();
      }
    }, 200);

    window.addEventListener('resize', this.redrawBoard.bind(this));
    window.addEventListener('load', this.redrawBoard.bind(this));
  }

  componentWillUpdate() {
    if (this.state.mounted) {
      this.redrawBoard();
    }
  }

  drawBoard(frame) {
    let board = frame.Board;
    let grid = ReactDom.findDOMNode(this.refs.gameGrid);
    let hexagonGrid = new HexagonGrid(grid, {
      radius: Math.min(
        1200 / board.Width / 2,
        700 / board.Height / Math.sqrt(3)
      )
    });

    hexagonGrid.drawHexGrid(
      this.generateBoardOptions(frame)
    );
  }

  generateBoardOptions(data) {
    let renderBoard = data.Board.Cells.reduce((board, row) => {
      row.forEach(cell => {
        let gridCell = this.props.emptyCell;
        if (cell.Filled) {
          gridCell = this.props.filledCell;
        }

        board['' + cell.X + cell.Y] = extend(gridCell);
      });
      return board;
    }, {});

    data.Unit.Members.forEach((cell) => {
      renderBoard['' + cell.X + cell.Y] = extend(this.props.droppingCell);
    });

    let pivot = data.Unit.Pivot;
    renderBoard['' + pivot.X + pivot.Y].Dot = true;

    return {
      columns: data.Board.Width,
      rows: data.Board.Height,
      board: renderBoard
    };
  }

  nextFrame(desiredFrame) {
    var frame = Math.min(
      desiredFrame || this.state.frameIndex + 1,
      this.props.gameData.length - 1
    );

    this.setState({
      frameIndex: frame
    });
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

  togglePause(paused) {
    if (this.state.paused === paused) return;

    this.setState({
      paused: paused
    });
  }

  render() {
    this.context.onSetTitle('ICFP!');
    let gameData = this.props.gameData[this.state.frameIndex] || {};

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
            <h5>Current Frame: {this.state.frameIndex}</h5>
            <h5>Current AI: {gameData.AI}</h5>
            <h5>Current Score: {gameData.Score}</h5>
            <br/>
            <button onClick={this.nextFrame.bind(this, this.state.frameIndex + 1)}>Next Frame</button>
            <button onClick={this.togglePause.bind(this, !this.state.paused)}>
              {this.state.paused ? 'Resume' : 'Pause'}
            </button>
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
    fill: '#fa0',
    stroke: '#000'
  },
  droppingCell: {
    fill: '#f10',
    stroke: '#000'
  },
  gridWidth: 1000,
  gridHeight: 700
};

export default MainPage;
