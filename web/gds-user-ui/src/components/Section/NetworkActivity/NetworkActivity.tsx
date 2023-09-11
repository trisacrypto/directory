import { VStack } from '@chakra-ui/react';
import { VictoryAxis, VictoryChart, VictoryLegend, VictoryLine, VictoryTheme } from 'victory';
import formatDisplayedDate from 'utils/formatDisplayedDate';
import useFetchNetworkActivity from './useFetchNetworkActivity';

const NetworkActivity = () => {
    const { data } = useFetchNetworkActivity();
    const mainnetData = data?.mainnet;
    const testnetData = data?.testnet;
    
    // The victory time scale requires dates to be in a Date object.
    mainnetData?.forEach((date: any) => {
      date.x = formatDisplayedDate(date.x);
    });
    
    testnetData?.forEach((date: any) => {
      date.x = formatDisplayedDate(date.x);
   });
   
   const legendData = [
    { name: "TestNet", symbol: { fill: "black" } },
    { name: "MainNet", symbol: { fill: "#F1511B" } }
  ];
  
  // Add padding to the axis labels to prevent overlap with the axis ticks.
  const sharedAxisStyles = {
    axisLabel: { padding: 35, fontWeight: 500, color: "black" },
  };

  return (
  <section>
    <VStack maxW={'5xl'} margin="auto" marginTop={10}>
      <VictoryChart
        domain={{ y: [0, 50] }}
        domainPadding={{ x: 1, y: 20 }}
        width={600} 
        style={{ background: { fill: "#F7F9FB" } }} 
        theme={VictoryTheme.material}
        scale={{ x: "time", y: "linear" }}
        >
          <VictoryLegend 
            data={legendData} 
            orientation="horizontal" 
            x={225}
            gutter={20}
            />
            <VictoryAxis fixLabelOverlap={true} />
            <VictoryAxis dependentAxis label="Network Activity" style={sharedAxisStyles} />
            <VictoryLine 
              data={testnetData} 
              x="date" 
              y="events"
              style={{ data: { stroke: "black" } }} 
              />
              <VictoryLine 
                data={mainnetData} 
                x="date" 
                y="events"
                style={{ data: { stroke: "#F1511B" } }}
                />
      </VictoryChart>
    </VStack>
  </section>
  );
};

export default NetworkActivity;
