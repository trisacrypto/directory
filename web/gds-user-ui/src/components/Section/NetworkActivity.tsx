import { VStack } from '@chakra-ui/react';
import { VictoryAxis, VictoryChart, VictoryLegend, VictoryLine, VictoryTheme, VictoryZoomContainer } from 'victory';
import { mockNetworkActivityData } from './_mocks_';
import formatDisplayedDate from 'utils/formatDisplayedDate';

const legendData = [
    { name: "TestNet", symbol: { fill: "black" } },
    { name: "MainNet", symbol: { fill: "#F1511B" } },
];

const mainnetData = mockNetworkActivityData?.networkActivity.mainnet;
const testnetData = mockNetworkActivityData?.networkActivity.testnet;

// The victory time scale requires dates to be in a Date object.
mainnetData?.forEach((d: any) => {
    d.x = formatDisplayedDate(d.x);
});

testnetData?.forEach((d: any) => {
    d.x = formatDisplayedDate(d.x);
});

// Add padding to the axis labels to prevent overlap with the axis ticks.
const sharedAxisStyles = {
    axisLabel: { padding: 35, fontWeight: 500, color: "black" },
};

const NetworkActivity = () => {
    return (
        <section>
          <VStack maxW={'5xl'} margin="auto" marginTop={10}>
            <VictoryChart 
              domainPadding={{ y: 20 }}
              containerComponent={<VictoryZoomContainer zoomDimension="x" />}
              width={600} 
              style={{ background: { fill: "#F7F9FB" } }} 
              theme={VictoryTheme.material}
              scale={{ x: "time" }} 
            >
              <VictoryLegend 
                data={legendData} 
                title="Network Activity |" 
                orientation="horizontal" 
                titleOrientation="left" 
                x={150} 
              />
              <VictoryAxis fixLabelOverlap={true} />
              <VictoryAxis dependentAxis label="Events" style={sharedAxisStyles} />
              <VictoryLine data={testnetData} style={{ data: { stroke: "black" } }} />
              <VictoryLine data={mainnetData} style={{ data: { stroke: "#F1511B" } }}/>
            </VictoryChart>
          </VStack>
        </section>
    );
};

export default NetworkActivity;
