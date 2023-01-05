import { Stack, chakra, ModalFooter, Button, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import CheckboxFormControl from 'components/ui/CheckboxFormControl';
import InputFormControl from 'components/ui/InputFormControl';
import { SubmitHandler, useFormContext } from 'react-hook-form';

type AddNewVaspFormProps = {
  onSubmit: SubmitHandler<any>;
  isCreatingVasp: boolean;
  closeModal: () => void;
};

function AddNewVaspForm({ onSubmit, isCreatingVasp, closeModal }: AddNewVaspFormProps) {
  const {
    handleSubmit,
    register,
    watch,
    formState: { errors }
  } = useFormContext();
  const accept = watch('accept');

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Text>
        <Trans>
          Please input the name of the new managed Virtual Asset Service Provider (VASP). When the
          entity is created, you will have the ability to add collaborators, start and complete the
          certificate registration process, and manage the VASP account. Please acknowledge below
          and provide the name of the entity.
        </Trans>
      </Text>

      <Stack pt={4}>
        <InputFormControl
          controlId="name"
          isInvalid={!!errors.name}
          data-testid="name"
          formHelperText={errors.name?.message as string}
          {...register('name')}
          label={
            <>
              <chakra.span fontWeight={700}>
                <Trans>VASP Name</Trans>
              </chakra.span>{' '}
              (<Trans>required</Trans>)
            </>
          }
        />
        <InputFormControl
          controlId="domain"
          isInvalid={!!errors.domain}
          data-testid="domain"
          formHelperText={errors.domain?.message as string}
          placeholder="https://"
          {...register('domain')}
          label={
            <>
              <chakra.span fontWeight={700}>
                <Trans>VASP Domain</Trans>
              </chakra.span>{' '}
              (<Trans>required</Trans>)
            </>
          }
        />
      </Stack>
      <CheckboxFormControl
        controlId="accept"
        data-testid="accept"
        {...register('accept', { required: true })}
        colorScheme="gray"
        borderColor="black">
        <Trans>
          TRISA is a network of trusted members. I acknowledge that the new VASP has a legitimate
          business purpose to join TRISA.
        </Trans>
      </CheckboxFormControl>
      <ModalFooter display="flex" flexDir="column" justifyContent="center" gap={2}>
        <Button
          isLoading={isCreatingVasp}
          bg="orange"
          _hover={{ bg: 'orange' }}
          type="submit"
          minW={150}
          isDisabled={!accept || isCreatingVasp}>
          <Trans id="Next">Create</Trans>
        </Button>
        <Button variant="ghost" onClick={closeModal} disabled={isCreatingVasp}>
          <Trans id="Cancel">Cancel</Trans>
        </Button>
      </ModalFooter>
    </form>
  );
}

export default AddNewVaspForm;
