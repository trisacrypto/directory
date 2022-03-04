import { DeleteIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Heading,
  HStack,
  Icon,
  Stack,
  Text,
  Tooltip,
  VStack,
} from "@chakra-ui/react";
import FormLayout from "layouts/FormLayout";
import LegalPersonForm from "./LegalPersonForm";

type LegalPersonProps = {};
const LegalPerson: React.FC<LegalPersonProps> = () => {
  return (
    <FormLayout>
      <Heading size="md">Addresses</Heading>
      <Text size="sm">Enter at least one geographic address</Text>
      {/* TODO: add Formik fieldarray to allow adding new address line */}
      <VStack width="100%" align="start" spacing={4}>
        <HStack width="100%" spacing={4}>
          <Box flex={1}>
            <LegalPersonForm />
          </Box>
          <Box alignSelf="flex-end" w={10} pb="25.1px">
            <Tooltip label="Delete the address line">
              <Button borderRadius={0}>
                <Icon as={DeleteIcon} />
              </Button>
            </Tooltip>
          </Box>
        </HStack>
        <Box>
          <Button borderRadius={0}>Add Address</Button>
        </Box>
      </VStack>
    </FormLayout>
  );
};

export default LegalPerson;
