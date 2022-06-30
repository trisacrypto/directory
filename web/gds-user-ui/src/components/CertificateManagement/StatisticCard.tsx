import { Heading, VStack } from "@chakra-ui/react";
import React from "react";

type StatisticCardProps = {
    title: string;
    total: number;
};


function StatisticCard({ title, total }: StatisticCardProps) {
    return <>
        <VStack border={"1px solid rgba(196, 196, 196, 0.58)"} p={5} flexGrow={1}>
            <Heading fontSize={'1.2rem'} textTransform={"capitalize"}>{title}</Heading>
            <Heading fontSize={'1.2rem'}>{total}</Heading>;
        </VStack>
    </>;
}

export default StatisticCard;
