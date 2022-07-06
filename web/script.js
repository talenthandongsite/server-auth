const INPUT_MAX = 200;
const PASSWORD_MIN = 4;
const KEYBOARD_EVENT_OFFSET_TIME = 50;

const accountInput = document.getElementById('accountInput');
const passwordInput = document.getElementById('passwordInput');
const submitButton = document.getElementById('submitButton');

const accountForm = new Form(accountInput, [
    FormValidator.nonEmptyString(),
    FormValidator.maxStringLength(INPUT_MAX)
]);

const passwordForm = new Form(passwordInput, [
    FormValidator.minStringLength(PASSWORD_MIN),
    FormValidator.maxStringLength(INPUT_MAX)
]);

accountInput.addEventListener('keydown', loginFormKeydown);
passwordInput.addEventListener('keydown', loginFormKeydown);
submitButton.addEventListener('click', signIn);

function loginFormKeydown(event) {
    setTimeout(() => {
        if (accountForm.prestine || passwordForm.prestine || !accountForm.checkValid() || !passwordForm.checkValid()) {
            submitButton.disabled = true;
            return;
        }
        submitButton.disabled = false;

        if (event.key == 'Enter') {
            signIn();
        }
    }, KEYBOARD_EVENT_OFFSET_TIME);
}

function signIn() {
    if (accountForm.prestine || passwordForm.prestine || !accountForm.checkValid() || !passwordForm.checkValid()) {
        return;
    }

    const account = accountForm.inputElement.value;
    const password = passwordForm.inputElement.value;

    Authentication.signIn(account, password);
}
