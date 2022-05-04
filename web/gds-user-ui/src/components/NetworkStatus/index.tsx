import { Stack, Text } from "@chakra-ui/react";

import { IoEllipse } from "react-icons/io5";

interface NetworkStatusProps {
  isOnline: boolean;
}
const NetworkStatus = (props: NetworkStatusProps) => {
  const { isOnline } = props;
  return (
    <Stack minHeight={82} bg={"white"} p={5} border="1px solid #C4C4C4">
      <Stack
        direction={"row"}
        justifyContent="space-between"
        alignItems="center"
      >
        <Text fontWeight={"bold"}> Network Status </Text>
        {isOnline ? (
          <IoEllipse fontSize="2rem" fill={"#34A853"} />
        ) : (
          <IoEllipse fontSize="2rem" fill={"#C4C4C4"} />
        )}
      </Stack>
    </Stack>
  );
};
NetworkStatus.defaultProps = {
  isOnline: true,
};
export default NetworkStatus;
