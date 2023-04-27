import { HStack, Button, FormControl } from '@chakra-ui/react';

import InputFormControl from 'components/ui/InputFormControl';
import { SubmitHandler, useForm } from 'react-hook-form';
import { t } from '@lingui/macro';
import { SearchIcon } from '@chakra-ui/icons';
import { Dispatch, SetStateAction } from 'react';
// import { useOrganizationListByName } from 'modules/dashboard/organization/useOrganizationListByName';
import { getOrganizationByName } from 'modules/dashboard/organization/organizationService';
type SearchVaspProps = {
  setSearchOrganization: Dispatch<SetStateAction<any>>;
};

const SearchVasp = (props: SearchVaspProps) => {
  console.log(props);

  // const { getAllOrganizations, organizations } = useOrganizationListByName();

  // const [showClose, setShowClose] = useState(false);
  const {
    handleSubmit,
    register,
    formState: { errors }
  } = useForm();

  const onSubmit: SubmitHandler<any> = (data) => {
    console.log('[data]', data);
    // fetch all organizations

    const d = getOrganizationByName(data.search as string);
    console.log('[data]', d);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <FormControl>
        <HStack>
          <InputFormControl
            controlId="search"
            size="md"
            width={'100%'}
            isInvalid={!!errors.search}
            placeholder={t`VASP Name`}
            data-testid="name"
            type="search"
            maxLength={200}
            formHelperText={errors.search?.message as string}
            {...register('search')}
          />
          <Button variant="outline" type="submit" spinnerPlacement="start">
            <SearchIcon />
          </Button>
        </HStack>
      </FormControl>
    </form>
  );
};

export default SearchVasp;
