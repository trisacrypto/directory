import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { fireEvent, render, screen } from 'utils/test-utils';
import InputFormControl from '.';

describe('<InputFormControl />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should ', () => {
    let formState: any = undefined;

    const Component = () => {
      const { register, formState: tempFormState } = useForm<{
        test: string;
      }>();
      formState = tempFormState;
      // eslint-disable-next-line no-unused-expressions
      formState.touchedFields;

      return (
        <div>
          <InputFormControl
            controlId="test"
            isInvalid
            {...register('test', { required: true })}
            formHelperText="Form Helper Text"
          />
        </div>
      );
    };

    const { debug } = render(<Component />);

    fireEvent.blur(screen.getByRole('textbox'), {
      target: {
        value: 'test'
      }
    });

    expect(formState.touchedFields.test).toBeDefined();
    expect(formState.isDirty).toBeFalsy();

    expect(screen.getByText(/Form Helper Text/i)).toBeInTheDocument();
    // debug();
  });
});
