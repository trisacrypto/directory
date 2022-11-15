import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import { ProfileBlock } from '..';

describe('<ProfileBlock />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render correctly props', () => {
    const title = 'This is title';
    const children = 'This is children';

    render(<ProfileBlock title={title}>{children}</ProfileBlock>);

    const titleEl = screen.getByTestId(/profile_block_title/i);
    expect(titleEl).toBeInTheDocument();
  });

  it('should match snapshot', () => {
    const title = 'This is title';
    const children = 'This is children';

    const { container } = render(<ProfileBlock title={title}>{children}</ProfileBlock>);

    expect(container).toMatchInlineSnapshot(`
      .emotion-0 {
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        -webkit-align-items: start;
        -webkit-box-align: start;
        -ms-flex-align: start;
        align-items: start;
        -webkit-flex-direction: column;
        -ms-flex-direction: column;
        flex-direction: column;
        width: 100%;
      }

      .emotion-0>*:not(style)~*:not(style) {
        margin-top: var(--chakra-space-5);
        -webkit-margin-end: 0px;
        margin-inline-end: 0px;
        margin-bottom: 0px;
        -webkit-margin-start: 0px;
        margin-inline-start: 0px;
      }

      .emotion-1 {
        font-family: var(--chakra-fonts-heading);
        font-weight: 700;
        font-size: var(--chakra-fontSizes-md);
        line-height: 1.2;
        text-transform: uppercase;
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        -webkit-column-gap: var(--chakra-space-4);
        column-gap: var(--chakra-space-4);
        -webkit-align-items: center;
        -webkit-box-align: center;
        -ms-flex-align: center;
        align-items: center;
      }

      .emotion-2 {
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        -webkit-align-items: start;
        -webkit-box-align: start;
        -ms-flex-align: start;
        align-items: start;
        -webkit-flex-direction: column;
        -ms-flex-direction: column;
        flex-direction: column;
        width: 100%;
      }

      .emotion-2>*:not(style)~*:not(style) {
        margin-top: var(--chakra-space-4);
        -webkit-margin-end: 0px;
        margin-inline-end: 0px;
        margin-bottom: 0px;
        -webkit-margin-start: 0px;
        margin-inline-start: 0px;
      }

      <div>
        <div
          class="chakra-stack emotion-0"
        >
          <h2
            class="chakra-heading emotion-1"
            data-testid="profile_block_title"
          >
            This is title
          </h2>
          <div
            class="chakra-stack emotion-2"
          />
        </div>
      </div>
    `);
  });
});
