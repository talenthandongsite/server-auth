const EMAIL_REGEX = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;

class FormValidator {

    static emailFormat() {
        return (value) => EMAIL_REGEX.test(value);
    }
    
    static maxStringLength(length) {
        if (length < 0) throw 'length cannot be a negative number';
        return (value) => {
            if (typeof value != 'string') return false;
            if (value.length > length) return false;
            return true;
        }
    }

    static minStringLength(length) {
        if (length < 0) throw 'length cannot be a negative number';
        return (value) => {
            if (typeof value != 'string') return false;
            if (value.length < length) return false;
            return true;
        }
    }

    static nonEmptyString() {
        return (value) => {
            if (typeof value != 'string') return false;
            if (value.length == 0) return false;
            return true;
        }
    }
}


class Form {

    constructor(inputElement, validators) {
        this.inputElement = inputElement;
        this.validators = validators;
        this.prestine = true;
        this.inputElement.addEventListener('keydown', (event) => {
            this.prestine = false;
        });
    }

    checkValid() {
        if (this.prestine) return false;
        return !this.validators.some(validator => {
            return !validator(this.inputElement.value);
        });
    }
}