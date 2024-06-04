import { InfoIcon, CheckCircleIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Link, Stack, Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';

import FormLayout from 'layouts/FormLayout';
import { ReactNode } from 'react';

type CertificateSectionProps = {
  step: number;
  title?: string;
  description?: string;
  isSaved?: boolean;
  isSubmitted?: boolean;
};

const getSection = (
  step: number
): { title: string; description?: string | ReactNode } | undefined => {
  switch (step) {
    case 1:
      return {
        title: t`Section 1: Basic Details`
      };
    case 2:
      return {
        title: t`Section 2: Legal Person`,
        description: t`Please enter the information that identifies your organization as a Legal Person. This form represents the ${(
          <Link href="https://intervasp.org/" isExternal color="link">
            IVMS 101
          </Link>
        )} data structure for legal persons and is strongly suggested for use as KYC (Know your Counterparty) or CDD (Customer Due Diligence) information exchanged in TRISA transfers.`
      };
    case 3:
      return {
        title: t`Section 3: Contacts`,
        description: t`Please supply contact information for representatives of your organization. All contacts will receive an email verification token and the contact emails must be verified before the registration can proceed.`
      };
    case 4:
      return {
        title: t`Section 4: TRISA Implementation`,
        description: t`Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please specify the details of your endpoint for certificate issuance.`
      };
    case 5:
      return {
        title: t`Section 5: TRIXO Questionnaire`,
        description: t`Please review the information provided, edit as needed, and submit to complete the registration form. After the information is reviewed, emails will be sent to the provided contacts for verification. Once verified, your TestNet certificate will be issued.`
      };
    case 6:
      return {
        title: t`Section 6: Review & Submit`,
        description: t`Please enter the information that identifies your organization as a Legal Person. This form represents the ${(
          <Link href="https://intervasp.org/" isExternal color="link">
            IVMS 101
          </Link>
        )} data structure for legal persons and is strongly suggested for use as KYC (Know your Counterparty) or CDD (Customer Due Diligence) information exchanged in TRISA transfers.`
      };
    case 7:
      return {
        title: t`Section 6: Review & Submit`,
        description: t`Your registration form has been successfully submitted. You will receive a confirmation email from admin@rotational.io. In the email, you will receive instructions on next steps. Return to your dashboard to monitor the status of your registration and certificate.`
      };
    default:
      return {
        title: t`Section 1: Basic Details`
      };
  }
};

const CertificateSection: React.FC<CertificateSectionProps> = ({
  step,
  isSaved,
  isSubmitted,
  title,
  description
}) => {
  return (
    <Stack>
      <HStack>
        <Heading size="md">{title || getSection(step)?.title}</Heading>
        <Box>
          {isSaved && step !== 7 && (
            <>
              <Icon as={InfoIcon} color="#F29C36" w={7} h={7} />
              <Text as={'span'} pl={2}>
                {' '}
                (<Trans id="Not Saved">Not Saved</Trans>)
              </Text>
            </>
          )}
          {!isSaved && step !== 7 && (
            <>
              <Icon as={CheckCircleIcon} color="#34A853" w={7} h={7} />
              <Text as={'span'} pl={2}>
                {' '}
                <Trans id="(Saved)">(Saved)</Trans>
              </Text>
            </>
          )}

          {step === 7 && isSubmitted && (
            <>
              <Icon as={CheckCircleIcon} color="#34A853" w={7} h={7} />
              <Text as={'span'} pl={2}>
                (<Trans id="Submitted">Submitted</Trans>)
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

CertificateSection.defaultProps = {
  step: 1
};

export default CertificateSection;
