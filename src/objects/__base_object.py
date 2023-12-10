class BaseObject(object):
    def __repr__(self):
        return str(self.__dict__)
