import { Box, Heading, Stack, Icon, HStack } from "@chakra-ui/react";
import { CheckCircleIcon } from "@chakra-ui/icons";

import BasicDetailsForm from "components/BasicDetailsForm";
import { FormikProps } from "formik";

type BasicDetailsProps = {
  formik: FormikProps<any>;
};

const BasicDetails: React.FC<BasicDetailsProps> = () => {
  return (
    <Stack
      spacing={5}
      paddingX="39px"
      paddingY="27px"
      border="3px solid #E5EDF1"
      borderRadius="md"
    >
      <HStack>
        <Heading size="md">Section 1: Basic Details</Heading>{" "}
        <Box>
          <Icon as={CheckCircleIcon} color="green.300" /> (saved)
        </Box>
      </HStack>
      <Box w={{ base: "100%", lg: "715px" }}>
        <BasicDetailsForm />
      </Box>
    </Stack>
  );
};

export default BasicDetails;
