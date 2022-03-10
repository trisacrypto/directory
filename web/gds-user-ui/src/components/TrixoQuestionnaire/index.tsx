import { InfoIcon } from "@chakra-ui/icons";
import { Box, Heading, HStack, Icon, Stack, Text } from "@chakra-ui/react";
import TrixoQuestionnaireForm from "components/TrixoQuestionnaireForm";
import FormLayout from "layouts/FormLayout";

const TrixoQuestionnaire: React.FC = () => {
  return (
    <Stack spacing={4}>
      <HStack>
        <Heading size="md">Section 5: TRIXO Questionnaire</Heading>
        <Box>
          <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved)
        </Box>
      </HStack>
      <FormLayout>
        <Text>
          This questionnaire is designed to help TRISA members understand the
          regulatory regime of your organization. The information provided will
          help ensure that required compliance information exchanges are
          conducted correctly and safely. All verified TRISA members will have
          access to this information.
        </Text>
      </FormLayout>
      <TrixoQuestionnaireForm />
    </Stack>
  );
};

export default TrixoQuestionnaire;
