import React from 'react';
import ReactDOM from 'react-dom';
import { createTheme, ThemeProvider } from '@material-ui/core/styles';
import App from './App';

const theme = createTheme();

ReactDOM.render(
  <ThemeProvider theme={theme}>
    <App />
  </ThemeProvider>,
  document.getElementById('app'),
);
