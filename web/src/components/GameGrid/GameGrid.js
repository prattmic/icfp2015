import React, { PropTypes } from 'react';
import styles from './GameGrid.css';
import withStyles from '../../decorators/withStyles';

@withStyles(styles)
class GameGrid {

  static propTypes = {
    path: PropTypes.string.isRequired,
    content: PropTypes.string.isRequired,
    title: PropTypes.string
  };

  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired
  };

  render() {
    this.context.onSetTitle(this.props.title);
    return (
      <div className="GameGrid">
        <div className="GameGrid-container">

          <div dangerouslySetInnerHTML={{__html: this.props.content || ''}} />
        </div>
      </div>
    );
  }

}

export default GameGrid;
