from marshmallow import Schema, EXCLUDE, post_load, ValidationError


class BaseSchema(Schema):
    __obj__ = None

    class Meta:
        unknown = EXCLUDE


    # @post_load
    # def make_object(self, data, **kwargs):
    #     return self.__obj__(**data)

    def handle_error(self, exc, data, **kwargs):
        raise ValidationError(exc.messages)
