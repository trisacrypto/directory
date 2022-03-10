import { InfoIcon, CheckCircleIcon } from "@chakra-ui/icons";
import { Box, Heading, HStack, Icon, Stack, Text } from "@chakra-ui/react";

import FormLayout from "layouts/FormLayout";

type CertificateSectionProps = {
  step: number;
  title?: string;
  description?: string;
  isSaved?: boolean;
  isSubmitted?: boolean;
};

const getSection = (
  step: number
): { title: string; description?: string } | undefined => {
  switch (step) {
    case 1:
      return {
        title: `Section 1: Basic Details`,
      };
    case 2:
      return {
        title: `Section 2: Legal Person`,
        description: `Please enter the information that identify your organization as a Legal Person. This form represents the IVMS 101 data structure for legal persons and is strongly suggested for use as KYC or CDD information exchanged in TRISA transfers.`,
      };
    case 3:
      return {
        title: `Section 3: Contacts`,
        description: `Please supply contact information for representatives of your organization. All contacts will receive an email verification token and the contact email must be verified before the registration can proceed.`,
      };
    case 4:
      return {
        title: `Section 4: TRISA Implementation`,
        description: `Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please specify the details of your endpoint for certificate issuance.`,
      };
    case 5:
      return {
        title: `Section 5: TRIXO Questionnaire`,
        description: `Please review the information provided, edit as needed, and submit to complete the registration form. After the information is reviewed, you will be contacted to verify details. Once verified, your TestNet certificate will be issued.`,
      };
    case 6:
      return {
        title: `Section 6: Review & Submit`,
        description: `Please enter the information that identify your organization as a Legal Person. This form represents the IVMS 101 data structure for legal persons and is strongly suggested for use as KYC or CDD information exchanged in TRISA transfers.`,
      };
    case 7:
      return {
        title: `Section 6: Review & Submit`,
        description: `Your registration form has been successfully submitted. You will receive a confirmation email from admin@trisa.io. In the email, you will receive instructions on next steps. Return to your dashboard to monitor the status of your registration and certificate.`,
      };
    default:
      return {
        title: `Section 1: Basic Details`,
      };
  }
};

const CertificateSection: React.FC<CertificateSectionProps> = ({
  step,
  isSaved,
  isSubmitted,
  title,
  description,
}) => {
  return (
    <Stack>
      <HStack>
        <Heading size="md">{title || getSection(step)?.title}</Heading>
        <Box>
          {isSaved && step !== 7 && (
            <>
              <Icon as={InfoIcon} color="#F29C36" w={7} h={7} />
              <Text as={"span"} pl={2}>
                {" "}
                (Not Saved)
              </Text>
            </>
          )}
          {!isSaved && step !== 7 && (
            <>
              <Icon as={CheckCircleIcon} color="#34A853" w={7} h={7} />
              <Text as={"span"} pl={2}>
                {" "}
                (Saved)
              </Text>
            </>
          )}

          {step === 7 && isSubmitted && (
            <>
              <Icon as={CheckCircleIcon} color="#34A853" w={7} h={7} />
              <Text as={"span"} pl={2}>
                {" "}
                (Submitted)
              </Text>
            </>
          )}
        </Box>
      </HStack>
      {(description || getSection(step)?.description) && (
        <FormLayout>
          <Text>{description || getSection(step)?.description}</Text>
        </FormLayout>
      )}
    </Stack>
  );
};

export default CertificateSection;
