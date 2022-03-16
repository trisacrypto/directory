import { Grid, GridItem, Heading, Text, VStack } from '@chakra-ui/react';
import Regulations from 'components/Regulations';
import SwitchFormControl from 'components/SwitchFormControl';
import FormButton from 'components/ui/FormButton';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import FormLayout from 'layouts/FormLayout';

const TrixoQuestionnaireForm: React.FC = () => {
  return (
    <FormLayout spacing={4}>
      <SelectFormControl
        controlId="primaryNationalJuridiction"
        label="Primary National Jurisdiction"
      />

      <InputFormControl
        controlId="nameOfPrimaryRegulator"
        label="Name of Primary Regulator"
        formHelperText="The name of primary regulator or supervisory authority for your national jurisdiction"
      />
      <VStack align="start">
        <Heading size="md">Other Jurisdictions</Heading>
        <Text>Please add any other regulatory jurisdictions your organization complies with.</Text>
      </VStack>
      <FormButton borderRadius={5}>Add Jurisdiction</FormButton>

      <SelectFormControl
        label="Is your organization permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?"
        controlId="financial_transfers_permitted"
      />

      <VStack align="start">
        <Heading size="md">CDD & Travel Rule Policies</Heading>

        <SelectFormControl
          label=" Does your organization have a programme that sets minimum AML, CFT,
        KYC/CDD and Sanctions standards per the requirements of the
        jurisdiction(s) regulatory regimes where it is
        licensed/approved/registered?"
          controlId="has_required_regulatory_program"
        />

        <VStack>
          <Text>
            Does your organization conduct KYC/CDD before permitting its customers to send/receive
            virtual asset transfers?
          </Text>
          <SwitchFormControl
            label="Conducts KYC before virtual asset transfers"
            controlId="conducts_customer_kyc"
          />
        </VStack>
      </VStack>

      <VStack align="start" w="100%">
        <Text>At what threshold and currency does your organization conduct KYC?</Text>
        <Grid templateColumns={{ base: '1fr 1fr', md: '2fr 1fr' }} gap={6} width="100%">
          <GridItem>
            <InputFormControl type="number" label="" controlId="country" />
          </GridItem>
          <GridItem>
            <SelectFormControl controlId="kyc_threshold_currency" />
          </GridItem>
        </Grid>
        <Text fontSize="sm" color="whiteAlpha.600" mt="0 !important">
          Threshold to conduct KYC before permitting the customer to send/receive virtual asset
          transfers
        </Text>
      </VStack>
      <VStack>
        <Text>
          Is your organization required to comply with the application of the Travel Rule standards
          in the jurisdiction(s) where it is licensed/approved/registered?
        </Text>
        <SwitchFormControl label="Must comply with Travel Rule" controlId="conducts_customer_kyc" />
      </VStack>
      <VStack align="start" w="100%">
        <Heading size="md">Applicable Regulations</Heading>
        <Text fontSize="sm">
          Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g.
          "FATF Recommendation 16"
        </Text>
        <Regulations />
      </VStack>

      <VStack align="start" w="100%" spacing={4}>
        <Heading size="md">Data Protection Policies</Heading>
        <VStack align="start" w="100%">
          <Text>Is your organization required by law to safeguard PII?</Text>
          <SwitchFormControl label="Must safeguard PII" controlId="must_safeguard_pii" />
        </VStack>
        <VStack align="start" w="100%">
          <Text>
            Does your organization secure and protect PII, including PII received from other VASPs
            under the Travel Rule?
          </Text>
          <SwitchFormControl label="Safeguards PII" controlId="safeguard_pii" />
        </VStack>
      </VStack>
    </FormLayout>
  );
};

export default TrixoQuestionnaireForm;
