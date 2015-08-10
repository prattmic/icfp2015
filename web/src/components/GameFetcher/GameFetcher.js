import React, { PropTypes } from 'react';
import ReactDom from 'react-dom';
import document from 'global/document';
import extend from 'xtend';
import MainPage from '../MainPage';
import styles from './GameFetcher.css';
import url from 'url';
import withStyles from '../../decorators/withStyles';
import window from 'global/window';
import xhr from 'xhr';

@withStyles(styles)
class GameFetcher extends React.Component {

  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired
  };

  state = {};

  componentDidMount() {
    this.fetchQualifiers();
    this.fetchNewGame();
  }

  fetchNewGame(options) {
    options = options || {};
    this.setState({
      gameData: null,
      qualifierName: options.qualifier || 'default'
    });

    xhr(url.format({
      pathname: '/api/newgame',
      query: {
        ai: options.ai || 'treeai',
        qualifier: options.qualifier
      }
    }), (err, resp, body) => {
      this.setState({
        gameData: JSON.parse(body)
      })
    });
  }

  fetchQualifiers() {
    xhr('/api/getqualifier', (err, resp, body) => {
      this.setState({
        qualifiers: JSON.parse(body)
      })
    });
  }

  render() {
    if (!this.state.gameData || !this.state.qualifiers) {
      return (
        <h1 className="loading-game">Game loading... Qualfier {this.state.qualifierName}</h1>
      );
    }

    return (
      <div className="GameFetcher">
        <MainPage
          fetchNewGame={this.fetchNewGame.bind(this)}
          board={this.state.gameData.Board}
          gameData={this.state.gameData.Frames}
          gridWidth={this.props.gridWidth}
          gridHeight={this.props.gridHeight}
          qualifiers={this.state.qualifiers}/>
      </div>
    );
  }
}

export default GameFetcher;
