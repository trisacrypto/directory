import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, waitFor } from 'utils/test-utils';
import Regulations from '.';

describe('<Regulations />', () => {
  beforeAll(async () => {
    await waitFor(() => {
      dynamicActivate('en');
    });
  });

  it('should should match snapshot', () => {
    const { container } = render(<Regulations name="" />);

    expect(container).toMatchInlineSnapshot(`
      <div>
        <div
          class="chakra-stack css-70kfd8"
        >
          <button
            class="chakra-button css-19d9t9j"
            type="button"
          >
            Add Regulation
          </button>
        </div>
      </div>
    `);
  });
});
