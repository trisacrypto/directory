import { VStack } from '@chakra-ui/react';
import { VictoryAxis, VictoryChart, VictoryLegend, VictoryLine, VictoryTheme } from 'victory';
import { mockNetworkActivityData } from './_mocks_';

const legendData = [
    { name: "TestNet", symbol: { fill: "black" } },
    { name: "MainNet", symbol: { fill: "#F1511B" } },
];

const mainnetData = mockNetworkActivityData?.networkActivity.mainnet;
const testnetData = mockNetworkActivityData?.networkActivity.testnet;

const NetworkActivity = () => {
    return (
        <>
        <VStack maxW={'5xl'}>
        <VictoryChart domainPadding={20} style={{ background: { fill: "#F7F9FB" } }} theme={VictoryTheme.material} width={600}>
            <VictoryLegend data={legendData} title="Network Activity" orientation="horizontal" titleOrientation="left" />
            <VictoryAxis tickValues={mainnetData.map((x: any) => x.x)} />
            <VictoryAxis dependentAxis />
            <VictoryLine data={testnetData} style={{ data: { stroke: "black" } }} />
            <VictoryLine data={mainnetData} style={{ data: { stroke: "#F1511B" } }}/>
        </VictoryChart>
        </VStack>
        </>
    );
};

export default NetworkActivity;
