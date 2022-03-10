import React, { FC } from "react";
import { Box, Text } from "@chakra-ui/react";

import EllipseIcon from "components/Icon/EllipseIcon";

interface NetworkStatusProps {
  isOnline: boolean;
}
const NetworkStatus = (props: NetworkStatusProps) => {
  const { isOnline } = props;
  return (
    <Box
      bg={"white"}
      minHeight={67}
      minWidth={246}
      pt={5}
      mt={10}
      mx={5}
      px={5}
      border="2px solid #C4C4C4"
      fontFamily={"Open Sans"}
    >
      <Box pb={2} display={"flex"} justifyContent={"space-between"}>
        <Text fontWeight={"bold"}> Network Status </Text>
        {isOnline ? (
          <EllipseIcon fill={"#34A853"} />
        ) : (
          <EllipseIcon fill={"#C4C4C4"} />
        )}
      </Box>
    </Box>
  );
};
NetworkStatus.defaultProps = {
  isOnline: true,
};
export default NetworkStatus;
