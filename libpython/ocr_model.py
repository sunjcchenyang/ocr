import json
def init():
    pass

def ocr_recog(img):
    r_list =[[0,1,'在水电费asdf',50,100,3,10,20,20,30,30,40]]
    return get_result(r_list)

def get_result(result_list):
    '''

    :param result_list: 包含检测与识别结果的列表，每一个元素是一个包含单条结果的列表
    :return: 返回相应的词典result['num_of_results'] 为结果数量，result[i]表示第i条结果
    例：result[i][typeID] 表示第i条结果的typeID
    '''
    results = []
    #results['num_of_results'] = len(result_list)
    for i, result_item in enumerate(result_list):

        result_dict = {}
        result_dict['TypeID'] = result_item[0]
        result_dict['NumID'] = result_item[1]
        result_dict['info'] = result_item[2]
        result_dict['Wide'] = result_item[3]
        result_dict['Hight'] = result_item[4]
        result_dict['CoordType'] = result_item[5]
        if result_item == 1:
            result_dict['absolutCoorX'] = result_item[6]
            result_dict['absolutCoorY'] = result_item[7]
            result_dict['relativeCoorX'] = ''
            result_dict['relativeCoorY'] = ''

        elif result_item == 2:
            result_dict['absolutCoorX'] = ''
            result_dict['absolutCoorY'] = ''
            result_dict['relativeCoorX'] = result_item[10]
            result_dict['relativeCoorY'] = result_item[11]

        elif result_item == 3:
            result_dict['absolutCoorX'] = result_item[6]
            result_dict['absolutCoorY'] = result_item[7]
            result_dict['relativeCoorX'] = result_item[10]
            result_dict['relativeCoorY'] = result_item[11]

        result_dict['relativeTypeID'] = result_item[8]
        result_dict['relativeNumID'] = result_item[9]
        result_json = json.dumps(result_dict)

        results.append(result_json)
    return results



#DEMO
# r_list = ocr_recog(img=None)
# print(r_list)
# print(type(r_list[0]))



