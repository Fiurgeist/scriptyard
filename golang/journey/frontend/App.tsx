import React from 'react';
import { ResponsiveLine, PointTooltipProps } from '@nivo/line';
import { BasicTooltip } from '@nivo/tooltip';
import { Paper } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(() => ({
  root: {
    height: '98vh',
    width: '98vh',
  },
}));

const randomHex = () => Math.floor(Math.random() * 200 + 20).toString(16);
const randomColor = () => `#${randomHex()}${randomHex()}${randomHex()}`;

const Tooltip: React.FC<PointTooltipProps> = ({ point }: PointTooltipProps) => {
  return (
    <BasicTooltip
      id={
        <span>
          <strong>Route {point.serieId}</strong> | x: <strong>{point.data.xFormatted}</strong>, y:
          <strong>{point.data.yFormatted}</strong>
        </span>
      }
      enableChip
      color={point.serieColor}
    />
  );
};

const App: React.FC = () => {
  const classes = useStyles();
  const [data, setData] = React.useState([]);

  React.useEffect(() => {
    const fetchRoutes = () =>
      fetch('http://localhost:8080/journeys', {
        method: 'GET',
        headers: {
          Accept: 'application/json; version=1',
        },
      })
        .then((response) => response.json())
        .then((parsed) => setData(parsed.journeys));

    fetchRoutes();
    const interval = setInterval(() => {
      fetchRoutes();
    }, 1000 * 15);
    return () => {
      clearInterval(interval);
    };
  }, []);

  return (
    <>
      <Paper className={classes.root}>
        <ResponsiveLine
          data={data}
          margin={{ top: 80, right: 80, bottom: 80, left: 80 }}
          xScale={{ type: 'linear', min: 0, max: 1023 }}
          yScale={{ type: 'linear', min: 0, max: 1023, stacked: false }}
          axisTop={{
            legend: 'Journeys',
            legendPosition: 'middle',
            legendOffset: -50,
            tickSize: 5,
            tickPadding: 5,
            tickRotation: 0,
          }}
          axisRight={{
            tickSize: 5,
            tickPadding: 5,
            tickRotation: 0,
          }}
          axisBottom={{
            tickSize: 5,
            tickPadding: 5,
            tickRotation: 0,
            legend: 'X position',
            legendOffset: 45,
            legendPosition: 'middle',
          }}
          axisLeft={{
            tickSize: 5,
            tickPadding: 5,
            tickRotation: 0,
            legend: 'Y position',
            legendOffset: -50,
            legendPosition: 'middle',
          }}
          enablePoints={false}
          useMesh
          colors={randomColor}
          tooltip={Tooltip}
          animate={false}
        />
      </Paper>
    </>
  );
};

export default App;
