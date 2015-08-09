import React, { PropTypes } from 'react';
import document from 'global/document';
import HexagonGrid from '../../lib/hexagon-grid';
import styles from './MainPage.css';
import withStyles from '../../decorators/withStyles';
import window from 'global/window';

@withStyles(styles)
class MainPage {

  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired
  };

  componentDidMount() {
    let hexagonGrid = new HexagonGrid(document.querySelector('.hexagon-game-grid'), {
      radius: 50
    });

    let board = {
      '00': {
        fill: '#f00',
        stroke: '#000'
      }
    };

    hexagonGrid.drawHexGrid({
      board: board,
      columns: 7,
      rows: 5
    });
  }

  render() {
    this.context.onSetTitle('ICFP!');
    var gridStyle = {
      position: 'relative'
    };

    return (
      <div className="MainPage">
        <div className="MainPage-container">
          <canvas className="hexagon-game-grid" width="1000" height="700">
          </canvas>
        </div>
      </div>
    );
  }
}

export default MainPage;
