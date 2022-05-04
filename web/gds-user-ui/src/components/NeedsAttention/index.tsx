import { Box, Text, Stack } from "@chakra-ui/react";
import FormButton from "components/ui/FormButton";

const NeedsAttention = () => {
  return (
    <Stack
      minHeight={67}
      bg={"#D8EAF6"}
      p={5}
      border="1px solid #DFE0EB"
      fontSize={18}
    >
      <Stack direction={"row"} spacing={3} alignItems="center">
        <Stack direction={["column", "row"]} spacing={3}>
          <Text fontWeight={"bold"}> Needs Attention </Text>
          <Text> Complete Testnet Registration </Text>
        </Stack>
        <Box>
          <FormButton width={142} borderRadius={5}>
            Start
          </FormButton>
        </Box>
      </Stack>
    </Stack>
  );
};

export default NeedsAttention;
