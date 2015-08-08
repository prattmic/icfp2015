import React, { PropTypes } from 'react';
import document from 'global/document';
import HexagonGrid from '../../lib/hexagon-grid';
import styles from './MainPage.css';
import withStyles from '../../decorators/withStyles';
import window from 'global/window';

@withStyles(styles)
class MainPage {

  static propTypes = {
  };

  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired
  };

  componentDidMount() {
    let hexagonGrid = new HexagonGrid(document.querySelector('.hexagon-game-grid'), 50);
    hexagonGrid.drawHexGrid(7, 10, 50, 50, true);
    //
    //let g;
    //function scan() {
    //  let opts = {
    //    element: gameGrid,
    //    width: 800,
    //    height: 500,
    //    spacing: 4
    //  };
    //
    //  let hexes = {
    //    width: 45, height: 45, n: 10
    //  };
    //
    //  g = grid(opts, hexes);
    //  console.log(g.grid);
    //}
    //
    //scan();
    //window.addEventListener('resize', scan);
    //window.addEventListener('load', scan);
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
