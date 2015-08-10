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

  state = {
    frameIndex: 0,
    interval: 200
  };

  componentDidMount() {
    this.setState({
      board: this.props.board,
      mounted: true
    });

    this.playFrames(this.state.interval);

    window.addEventListener('resize', this.redrawBoard.bind(this));
    window.addEventListener('load', this.redrawBoard.bind(this));
  }

  componentWillUpdate(nextProps, nextState) {
    if (!this.state.mounted) {
      return;
    }

    if (this.state.frameIndex !== nextState.frameIndex) {
      this.redrawBoard(nextState.frameIndex);
    }

    if (this.state.interval !== nextState.interval) {
      this.stopFrames();
      this.playFrames(nextState.interval);
    }
  }

  componentWillUnmount() {
    this.stopFrames();
  }

  drawBoard(frame) {
    let grid = ReactDom.findDOMNode(this.refs.gameGrid);
    let hexagonGrid = new HexagonGrid(grid, {
      radius: Math.min(
        1200 / this.state.board.Width / 2,
        700 / this.state.board.Height / Math.sqrt(3)
      )
    });

    let deltas = frame.BoardDelta;

    let updatedBoard = extend(this.state.board);
    deltas.forEach(delta => {
      updatedBoard.Cells[delta.X][delta.Y] = delta;
    });

    this.setState({
      board: updatedBoard
    });

    hexagonGrid.drawHexGrid(
      this.generateBoardOptions(updatedBoard, frame)
    );
  }

  generateBoardOptions(b, data) {
    let renderBoard = b.Cells.reduce((board, row) => {
      row.forEach(cell => {
        let gridCell = this.props.emptyCell;
        if (cell.Filled) {
          gridCell = this.props.filledCell;
        }

        board['x' + cell.X + 'y' + cell.Y] = extend(gridCell);
      });
      return board;
    }, {});

    if (data.Unit.Members) {
      data.Unit.Members.forEach((cell) => {
        renderBoard['x' + cell.X + 'y' +  cell.Y] = extend(this.props.droppingCell);
      });
    }

    let pivot = data.Unit.Pivot;
    if (pivot && renderBoard['x' + pivot.X + 'y' +  pivot.Y]) {
      renderBoard['x' + pivot.X + 'y' +  pivot.Y].Dot = true;
    }

    return {
      columns: this.state.board.Width,
      rows: this.state.board.Height,
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

  playFrames(interval) {
    this.setState({
      playingFrames: setInterval(() => {
        if (this.playNext()) {
          this.nextFrame();
        }
      }, interval)
    });
  }

  stopFrames() {
    if (this.state.playingFrames) {
      clearInterval(this.state.playingFrames);
    }
  }

  playNext() {
    if (this.state.paused) {
      return false;
    }

    return this.state.frameIndex < this.props.gameData.length;
  }

  redrawBoard(frameIndex) {
    this.drawBoard(this.props.gameData[frameIndex]);
  }

  setGameSpeed(e) {
    e.preventDefault();
    this.setState({
      interval: parseInt(this.refs.setSpeedInput.value.trim())
    });
  }

  submitNewGame(e) {
    e.preventDefault();
    this.props.fetchNewGame({
      ai: this.refs.selectAi.value.trim(),
      qualifier: this.refs.selectQualifier.value.trim(),
      repeater: this.refs.repeater.value.trim()
    });
  }

  togglePause(paused) {
    if (this.state.paused === paused) return;

    this.setState({
      paused: paused
    });

    if (paused) {
      this.stopFrames();
    } else {
      this.playFrames(this.state.interval);
    }
  }

  render() {
    this.context.onSetTitle('ICFP!');
    let gameData = this.props.gameData[this.state.frameIndex] || {};

    var qualifiers = this.props.qualifiers.map((qualifier) => {
      return (
        <option key={qualifier}>
          {qualifier}
        </option>
      );
    });

    var aiList = ['repeaterai', 'treeai', 'simpleai', 'lookaheadai'].map(ai => {
      return (
        <option key={ai}>
          {ai}
        </option>
      );
    });

    let interval = this.state.interval;

    return (
      <div className="MainPage">
        <div className="MainPage-container">
          <div className="hexagon-container">
            <canvas className="hexagon-game-grid" ref="gameGrid"
                    width={this.props.gridWidth} height={this.props.gridHeight}>
            </canvas>
          </div>
          <div className="panel-container">
            <div className="panel-items">
              <h4 className="max-width">Game Controls</h4>
              <span className="max-width">Current Frame: {this.state.frameIndex}</span>
              <br/>
              <span className="max-width">Current AI: {gameData.AI}</span>
              <br/>
              <span className="max-width">Current Score: {gameData.Score}</span>
              <br/>
              <button onClick={this.nextFrame.bind(this, this.state.frameIndex + 1)}>Next Frame</button>
              <button onClick={this.togglePause.bind(this, !this.state.paused)}>
                {this.state.paused ? 'Resume' : 'Pause'}
              </button>
              <br/>

              <form className="update-speed" onSubmit={this.setGameSpeed.bind(this)}>
                <label className="setSpeed-label max-width" htmlFor="setSpeed">Speed (ms)</label>
                <input className="setSpeed-select max-width" type="text"
                       name="setSpeed" ref="setSpeedInput" defaultValue={interval} />
                <input className="setSpeed-submit max-width" type="submit" value="Set Speed" />
              </form>

              <h4 className="max-width">New Game Options</h4>
              <form className="new-game" onSubmit={this.submitNewGame.bind(this)}>
                <label className="qualifier-label max-width" htmlFor="qualifier">Qualifier</label>
                <select className="qualifier-select max-width" name="qualifier" ref="selectQualifier">
                  {qualifiers}
                </select>
                <select className="ai-select max-width" name="ai" ref="selectAi">
                  {aiList}
                </select>
                <input className="max-width" type="text" name="repeater" ref="repeater"/>
                <input className="new-game-submit max-width" type="submit" value="New Game" />
              </form>
            </div>
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
