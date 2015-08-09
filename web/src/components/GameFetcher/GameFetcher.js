import React, { PropTypes } from 'react';
import ReactDom from 'react-dom';
import document from 'global/document';
import extend from 'xtend';
import MainPage from '../MainPage';
import styles from './GameFetcher.css';
import withStyles from '../../decorators/withStyles';
import window from 'global/window';
import xhr from 'xhr';

@withStyles(styles)
class GameFetcher extends React.Component {

  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired
  };

  state = {
  };

  componentDidMount() {
    xhr('/api/newgame', (err, resp, body) => {
      this.setState({
        gameData: JSON.parse(body)
      })
    });
  }

  render() {
    if (!this.state.gameData) {
      return (
        <div className="loading-game">Game loading...</div>
      );
    }

    return (
      <div className="GameFetcher">
        <MainPage gameData={this.state.gameData}
          gridWidth={this.props.gridWidth}
          gridHeight={this.props.gridHeight} />
      </div>
    );
  }
}

export default GameFetcher;
